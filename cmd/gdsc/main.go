package main

import (
	"os"
	"regexp"

	"github.com/alexflint/go-arg"
	"github.com/herlon214/gdsc/pkg/docker"
	"github.com/herlon214/gdsc/pkg/logger"
)

// Args contain th received args when execute 'gdsc'
type Args struct {
	From   string `arg:"required" help:"service that will be cloned if there is no service with the given --name"`
	Name   string `arg:"required" help:"service name that will be deployed"`
	Image  string `arg:"required" help:"new docker image url"`
	Domain string `arg:"required" help:"root domain to be used in the traefik host, eg: mycompany.org"`
}

// ParseArgs parse the received args
func ParseArgs() Args {
	var args Args
	arg.MustParse(&args)
	r := regexp.MustCompile("\\W")

	// Filter the new name to only words
	args.Name = r.ReplaceAllString(args.Name, "_")

	return args
}

func main() {
	// Logging

	var api = docker.Api{ApiUrl: "http://127.0.0.1:2375"}
	var log = logger.DefaultLogger()
	args := ParseArgs()

	// Try to get a crated service
	service := api.GetService(args.Name)

	// Check if the service is already created
	if service.Spec.Name == "" {
		log.Warningf("Service was %s not created yet, creating a new one based on %s ...", args.Name, args.From)
		service = api.GetService(args.From)

		newService := service

		// Change the service informations
		newService.Spec.TaskTemplate.ContainerSpec.Image = args.Image
		newService.Spec.Name = args.Name
		newService.Spec.Labels["traefik.frontend.rule"] = "Host: " + args.Name + "." + args.Domain

		response := api.CreateService(newService.Spec)
		log.Debugf("Service created with ID: %s", response.ID)
	} else { // Update a created service
		log.Debugf("Updating a existent service with name %s ...", args.Name)
		newService := service

		// Change the service informations
		newService.Spec.TaskTemplate.ContainerSpec.Image = args.Image
		newService.Spec.Name = args.Name
		newService.Spec.Labels["traefik.frontend.rule"] = "Host: " + args.Name + "." + args.Domain

		if api.UpdateService(*newService) == true {
			log.Noticef("Service %s updated successfully!", args.Name)
		} else {
			log.Errorf("Failure when updating service %s!", args.Name)

			// Exit with error status
			os.Exit(1)
		}
	}

}
