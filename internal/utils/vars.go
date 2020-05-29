package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

//list of key value pairs
type Vars map[string]string

//keys

const (
	PfxPath     string = "PFX_PATH"
	PfxPassword string = "PFX_PASSWORD"
	ServerName  string = "SERVER_NAME"

	Edge    string = "EDGE"
	Catalog string = "CATALOG"
	Userdn  string = "USER_DN"

	JsonConfigPath  string = "JSON_CONFIG_PATH"
	CatalogFileName string = "CATALOG_FILE_NAME"
)

func NewVars() (v Vars) {
	viper.SetConfigName("config")
	viper.AddConfigPath("..")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/Users/kyleg/Documents/new-scripts/")
	viper.AutomaticEnv()

	//default values

	v = Vars{
		PfxPath:     "/Users/kyleg/Documents/new-scripts/certs/quickstart.crt",
		PfxPassword: "",
		ServerName:  "",

		Edge:    "https://localhost:30000",
		Catalog: "/services/catalog/latest",
		Userdn:  "CN=quickstart,OU=Engineering,O=Decipher Technology Studios,L=Alexandria,ST=Virginia,C=US",

		JsonConfigPath:  "/Users/kyleg/Documents/new-scripts/json_config",
		CatalogFileName: "06.catalog.json",
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("Did Not find config file")
		} else {
			// Config file was found but another error was produced
			fmt.Print("Config file found but another error occurred")
			return
		}
	}

	for k := range v {
		//if viper has a value and its not blank ovewrite the default
		//add in the output for input vars here
		if envValue := viper.GetString(k); envValue != "" {
			v[k] = envValue
		}
	}

	return
}
