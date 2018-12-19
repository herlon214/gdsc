package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/herlon214/gdsc/pkg/docker"
	"github.com/herlon214/gdsc/pkg/logger"
	"github.com/spf13/cobra"
)

var ExportFrom string
var ExportTo string

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export a service to a json file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Logger
		log := logger.DefaultLogger()

		// Try to get a crated service
		api := docker.Api{ApiUrl: ApiUrl}
		service := api.GetService(ExportFrom)
		export, err := json.Marshal(service.Spec)

		if err != nil {
			panic(fmt.Errorf("Command failed with: %+v", err))
		}

		err = ioutil.WriteFile(ExportTo, []byte(export), 0644)

		if err != nil {
			panic(fmt.Errorf("Command failed with: %+v", err))
		} else {
			log.Noticef("Service %s exported to %s", ExportFrom, ExportTo)
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVar(&ExportFrom, "from", "", "Service name that will export")
	exportCmd.Flags().StringVar(&ExportTo, "to", "", "Output path to save the service json")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
