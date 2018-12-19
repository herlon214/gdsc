package cmd

import (
	"os"

	"github.com/herlon214/gdsc/pkg/docker"
	"github.com/herlon214/gdsc/pkg/logger"
	"github.com/spf13/cobra"
)

var CopyFrom string
var Name string
var Image string
var Domain string
var Auth string
var Daemon bool

// upsertCmd represents the upsert command
var upsertCmd = &cobra.Command{
	Use:   "upsert",
	Short: "Copy a service and create a new one overriding image and name",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Try to get a crated service
		api := docker.Api{ApiUrl: ApiUrl}
		service := api.GetService(Name)

		// Check if the service is already created
		if service.Spec.Name == "" {
			CreateService(api)
		} else { // Update a created service
			UpdateService(api, *service)
		}
	},
}

func init() {
	rootCmd.AddCommand(upsertCmd)

	upsertCmd.Flags().StringVar(&CopyFrom, "copy-from", "", "Service name that will be copied")
	upsertCmd.Flags().StringVar(&Name, "name", "", "New service name")
	upsertCmd.Flags().StringVar(&Image, "image", "", "Image that new service will use")
	upsertCmd.Flags().StringVar(&Domain, "domain", "", "Domain host to be set in traefik labels")
	upsertCmd.Flags().StringVar(&Auth, "auth", "", "Registry auth token")
	upsertCmd.Flags().BoolVarP(&Daemon, "daemon", "d", true, "Update service using docker daemon")

	upsertCmd.MarkFlagRequired("copy-from")
	upsertCmd.MarkFlagRequired("name")
}

// CreateService creates a new docker service
func CreateService(api docker.Api) {
	log := logger.DefaultLogger()
	log.Warningf("Service was %s not created yet, creating a new one based on %s ...", Name, CopyFrom)
	service := api.GetService(CopyFrom)
	headers := map[string]string{}

	// Check for auth
	if Auth != "" {
		headers["X-Registry-Auth"] = Auth
	}

	// Change the service informations
	service.Spec.TaskTemplate.ContainerSpec.Image = Image
	service.Spec.Name = Name

	// Check if must set a Traefik rule
	if Domain != "" {
		service.Spec.Labels["traefik.frontend.rule"] = "Host: " + Name + "." + Domain
	}

	response := api.CreateService(service.Spec, headers)
	log.Debugf("Service created with ID: %s", response.ID)

	// Update service with daemon (becuase of --with-registry not working)
	success := api.UpdateWithDaemon(*service)

	if success == true {
		log.Noticef("Service %s updated with daemon successfully!", Name)
	} else {
		log.Errorf("Failure when updating service %s!", Name)

		// Exit with error status
		os.Exit(1)
	}
}

// UpdateService updates a service with new docker image
func UpdateService(api docker.Api, service docker.Service) {
	log := logger.DefaultLogger()
	headers := map[string]string{}
	log.Debugf("Updating a existent service with name %s ...", Name)

	// Check for auth
	if Auth != "" {
		headers["X-Registry-Auth"] = Auth
	}

	// Change the service informations
	service.Spec.TaskTemplate.ContainerSpec.Image = Image
	var success bool

	// Check if must run as daemon
	if Daemon {
		success = api.UpdateWithDaemon(service)
	} else {
		success = api.UpdateService(service, headers)
	}

	if success == true {
		log.Noticef("Service %s updated with daemon successfully!", Name)
	} else {
		log.Errorf("Failure when updating service %s!", Name)

		// Exit with error status
		os.Exit(1)
	}
}
