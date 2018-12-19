package cmd

import (
	"fmt"
	"io/ioutil"

	"encoding/json"

	"github.com/herlon214/gdsc/pkg/docker"
	"github.com/herlon214/gdsc/pkg/logger"
	"github.com/spf13/cobra"
)

var ImportFrom string
var ImportName string

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a JSON file and create a service using it",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Logger
		log := logger.DefaultLogger()

		// Read file content
		file, err := ioutil.ReadFile(ImportFrom)

		if err != nil {
			panic(fmt.Errorf("Command failed with: %+v", err))
		}

		var service docker.Spec
		json.Unmarshal([]byte(file), &service)

		service.Name = ImportName

		// Try to create the service
		api := docker.Api{ApiUrl: ApiUrl}
		response := api.CreateRawService(service)

		if response.ID != "" {
			log.Noticef("Service %s created sucessfully!", ImportName)
		} else {
			log.Errorf("Failed to create service %s: %s", ImportName, response.Message)
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringVar(&ImportFrom, "from", "", "JSON file path to import")
	importCmd.Flags().StringVar(&ImportName, "name", "", "New service name")
}
