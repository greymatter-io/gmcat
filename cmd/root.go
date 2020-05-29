package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"local.only/gmcat/internal/utils"
)

//for use in sub commands
var directory string

var cfgFile string

var (
	rootCmd = &cobra.Command{
		Use:   "gmcat",
		Short: "Grey Matter Mesh Tool",
		Long:  `Grey Matter Mesh Tool is a tool to dynamically work with a deployed service mesh.  Currently supports catalog interactions.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("gmcat is a Grey Matter Catalog cli.  For use run gmcat --help")
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $PWD/config.yaml")
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initConfig() {

	utils.NewVars()

}
