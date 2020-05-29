package mesh

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/crypto/pkcs12"
)

// NewTLSClient takes a relative path to a pfx certificate, a password for that certificate, and the servername that validates that certificate. Then it converts the pfx and password to an X509 key pair, and returns an http client with TLS configured.
func newTLSClient(pfxPath string, pfxPassword string, serverName string) (client *http.Client, err error) {
	fmt.Println("Attempting to create a new tls client with pfx", pfxPath, "and servername", serverName)
	fmt.Println("Reading in pfx file...")
	pfxData, err := ioutil.ReadFile(pfxPath)
	if err != nil {
		fmt.Println("ERROR! reading pfx file:", err)
		return
	}
	fmt.Println("Converting pfx file data to pem blocks...")
	blocks, err := pkcs12.ToPEM(pfxData, pfxPassword)
	if err != nil {
		fmt.Println("ERROR! making pem blocks", err)
		return
	}
	fmt.Println("Converting pem blocks to x509 key pair...")
	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}
	cert, err := tls.X509KeyPair(pemData, pemData)
	if err != nil {
		fmt.Println("ERROR! making x509 key pair:", err)
		return
	}
	certs := x509.NewCertPool()
	certs.AppendCertsFromPEM(pemData)
	fmt.Println("Initializing tls client...")
	tlsClientConfig := &tls.Config{
		RootCAs:            certs,
		Certificates:       []tls.Certificate{cert},
		ServerName:         serverName,
		InsecureSkipVerify: true,
		ClientAuth:         tls.RequireAndVerifyClientCert,
	}
	transport := &http.Transport{TLSClientConfig: tlsClientConfig}
	client = &http.Client{
		Transport: transport,
	}
	fmt.Println("SUCCESS! created new tls client!")
	return
}
