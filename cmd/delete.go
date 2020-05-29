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
var addDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete catalog entries",
	Run: func(cmd *cobra.Command, args []string) {
		Delete()
	},
}

func init() {
	rootCmd.AddCommand(addDelete)
	addDelete.Flags().StringVarP(&directory, "file", "f", "", "-f <directory name>: pass the directory the catalog defn is located in")
}

func Delete() {
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

	var clusterToDelete string
	var zone = ""
	if len(directory) > 0 {
		pathName := jsonPath + "/" + directory + "/" + catalogFileName
		inputCluster, _ := catalog.FileToEntryCluster(pathName)
		clusterToDelete = inputCluster.ClusterName
		zone = inputCluster.ZoneName
	} else {
		list := mh.ListCatalogEntries()
		fmt.Print(utils.VSlice(list))

		clusterToDelete, _ = utils.InputFromList(list)
		zone = mh.GetCatalogEntry(clusterToDelete).ZoneName
		fmt.Printf("BBBBIIIIINNNNGGGGOOOO: %s\n", zone)
	}
	result := mh.DeleteCatalogEntry(clusterToDelete, zone)
	if result {
		fmt.Printf("Deleted %s entry from catalog\n\n", clusterToDelete)
	} else {
		fmt.Printf("There was an issue deleting %s entry from catalog\n\n", clusterToDelete)
	}
}
