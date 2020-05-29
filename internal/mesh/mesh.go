package mesh

import (
	"net/http"

	"local.only/gmcat/internal/utils"
)

type Handler struct {
	client *http.Client
	config *Config
}

type Config struct {
	// host
	Edge       string
	Catalog    string
	Controlapi string
	Userdn     string
}

func NewHandler(pfxPath, pfxPassword, serverName string, config *Config) (h *Handler) {
	client, err := newTLSClient(pfxPath, pfxPassword, serverName)
	utils.Check(err)

	h = &Handler{
		client: client,
		config: config,
	}
	return
}
