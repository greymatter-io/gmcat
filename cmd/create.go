package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"local.only/gmcat/internal/catalog"
	"local.only/gmcat/internal/mesh"
	"local.only/gmcat/internal/utils"
)

//cobra logic
var addCreate = &cobra.Command{
	Use:   "create",
	Short: "Create catalog entries",
	Run: func(cmd *cobra.Command, args []string) {
		Create()
	},
}

func init() {
	rootCmd.AddCommand(addCreate)
	addCreate.Flags().StringVarP(&directory, "file", "f", "", "-f <directory name>: pass the directory the catalog defn is located in")
}

func Create() {
	jsonPath := viper.Get("json_config_path").(string)
	catalogFileName := viper.Get("catalog_file_name").(string)
	pfxPath := viper.Get("pfx_path").(string)
	pfxPassword := viper.Get("pfx_password").(string)
	serverName := viper.Get("server_name").(string)
	config := &mesh.Config{
		Edge:    viper.Get("edge").(string),
		Catalog: viper.Get("catalog").(string),
		Userdn:  viper.Get("user_dn").(string),
	}

	mh := mesh.NewHandler(pfxPath, pfxPassword, serverName, config)

	if len(directory) > 0 {
		pathName := jsonPath + "/" + directory + "/" + catalogFileName
		inputCluster, err := catalog.FileToEntryCluster(pathName)
		utils.Check(err)
		inputCluster.PrintCluster()

		mh.CreateCatalogEntry(inputCluster)

	} else {
		fmt.Println("You need to specify a directory for the catalog configs with -f <directory>")
	}

}
