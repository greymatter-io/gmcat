package cmd

import (
	"fmt"

	"github.com/r3labs/diff"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"local.only/gmcat/internal/catalog"
	"local.only/gmcat/internal/mesh"
	"local.only/gmcat/internal/utils"
)

//cobra logic
var addSearch = &cobra.Command{
	Use:   "search",
	Short: "Search catalog entries",
	Run: func(cmd *cobra.Command, args []string) {
		Search()
	},
}

func init() {
	rootCmd.AddCommand(addSearch)
	addSearch.Flags().StringVarP(&directory, "file", "f", "", "-f <directory name>: pass the directory the catalog defn is located in")
}

//business logic
func Search() {
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

	catalog_entries := mh.ListCatalogEntries()
	fmt.Printf("\nclusterName(s):\n")
	fmt.Print(utils.VSlice(catalog_entries))

	if len(directory) > 0 {
		// parse input file to get the cluster from a file then check if that is in the mesh
		pathName := jsonPath + "/" + directory + "/" + catalogFileName
		inputCluster, err := catalog.FileToEntryCluster(pathName)
		utils.Check(err)

		clusterName := inputCluster.ClusterName

		fmt.Printf("Looking up clusterName: %s in mesh\n\n", clusterName)
		catalog_entry := mh.GetCatalogEntry(clusterName)
		_, ispresent := utils.Verify(clusterName, catalog_entries)

		if ispresent == true {
			fmt.Printf("\n-- Current entry in catalog: --\n")
			catalog_entry.PrintCluster()

			//diff from local entry to in-mesh entry
			changelog, _ := diff.Diff(inputCluster, catalog_entry)
			fmt.Printf("The change from local to in-mesh is:\n%s\n\n", changelog)

			catalog.DiffCatalogCluster(inputCluster, catalog_entry)

		} else {
			fmt.Printf("ClusterName (catalog entry) [%s] defined in %s is not currently in mesh", clusterName, pathName)
		}

	} else {
		//interactive mode

		var keepgoing = true
		for keepgoing {

			fmt.Println("Enter clusterName of the catalog entry you would like to see (from the list above): ")
			input, ispresent := utils.InputFromList(catalog_entries)
			catalog_entry := mh.GetCatalogEntry(input)
			catalog_entry.PrintCluster()
			if ispresent {
				keepgoing = false
			}
		}
	}
}
