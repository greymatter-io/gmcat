package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"local.only/gmcat/internal/catalog"
	"local.only/gmcat/internal/mesh"
	"local.only/gmcat/internal/utils"
)

//cobra logic
var addUpdate = &cobra.Command{
	Use:   "update",
	Short: "Update catalog entries",
	Run: func(cmd *cobra.Command, args []string) {
		Update()
	},
}

func init() {
	rootCmd.AddCommand(addUpdate)
	addUpdate.Flags().StringVarP(&directory, "file", "f", "", "-f <directory name>: pass the directory the catalog defn is located in")
}

func Update() {
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

	//create mesh handler
	mh := mesh.NewHandler(pfxPath, pfxPassword, serverName, config)

	//Get list of catalog entries and print them
	catalog_entries := mh.ListCatalogEntries()
	fmt.Print(utils.VSlice(catalog_entries))

	//chose the entry to edit, create file and backup
	var outputPath string
	var outputPathFile string
	var clusterToUpdate string

	var CatalogClusterMap = make(catalog.CatalogEntryMap)

	if len(directory) > 0 {
		// parse input file to get the cluster from a file then check if that is in the mesh
		pathName := jsonPath + "/" + directory + "/" + catalogFileName
		inputCluster, err := catalog.FileToEntryCluster(pathName)
		utils.Check(err)

		//see if the cluster  i want to update from is in the mesh
		clusterNameToCheck := inputCluster.ClusterName

		_, ispresent := utils.Verify(clusterNameToCheck, catalog_entries)
		if ispresent {
			fmt.Println(ispresent)
			//set the output path directory since it is present
			outputPath = fmt.Sprintf("%s/%s", jsonPath, directory)
			outputPathFile = fmt.Sprintf("%s/%s", outputPath, catalogFileName)

			CatalogClusterMap["current"] = inputCluster

			//backup the file
			backupOutputName := fmt.Sprint(catalogFileName, "-bk")
			_, _ = inputCluster.ClusterToFile(outputPath, backupOutputName)
		} else {
			panic("clusterName is not present in the mesh. cannot update what is not there.")
		}

		fmt.Printf("RIGHT HERE the outputPath is: %s\n", outputPath)

	} else {
		// get the entry and back it up to the proper directory
		var keepgoing = true
		for keepgoing {
			fmt.Println("Enter catalogName you would like edit (from the list above): ")

			var ispresent bool
			clusterToUpdate, ispresent = utils.InputFromList(catalog_entries)

			//set outputPath name based on the entry used
			outputPath = fmt.Sprintf("%s/%s", jsonPath, clusterToUpdate)

			//Create a file to represent what is in the mesh
			catalog_entry := mh.GetCatalogEntry(clusterToUpdate)

			//add the current entry from the mesh into a map and label it current
			CatalogClusterMap["current"] = catalog_entry

			//Create a file to modify and a backup to edit and process
			outputPathFile, _ = CatalogClusterMap["current"].ClusterToFile(outputPath, catalogFileName)
			backupOutputName := fmt.Sprint(catalogFileName, "-bk")
			_, _ = CatalogClusterMap["current"].ClusterToFile(outputPath, backupOutputName)

			if ispresent {
				keepgoing = false
			}
		}
	}

	//edit entry
	cmd := exec.Command("vim", outputPathFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	utils.Check(err)

	CatalogClusterMap["updated"], err = catalog.FileToEntryCluster(outputPathFile)
	utils.Check(err)

	//present the old and new entries as well as a diff of the two
	catalog.DiffCatalogCluster(CatalogClusterMap["current"], CatalogClusterMap["updated"])

	//remove the old entry and apply the new entry
	mh.DeleteCatalogEntry(CatalogClusterMap["current"].ClusterName, CatalogClusterMap["current"].ZoneName)

	//create new entry
	mh.CreateCatalogEntry(CatalogClusterMap["updated"])

}
