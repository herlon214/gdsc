package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/herlon214/gdsc/pkg/docker"
	"github.com/herlon214/gdsc/pkg/logger"
)

type Args struct {
	ImageURL        string
	BranchName      string
	ServiceCopyName string
	NewServiceName  string
	ServiceVersion  int
}

// ParseArgs parse the received args
func ParseArgs() Args {
	args := os.Args[1:]
	r := regexp.MustCompile("\\W")
	var BranchName = r.ReplaceAllString(args[1], "_")
	var NewServiceName = strings.Replace(args[0], "develop", BranchName, -1)
	var ImageURL = args[2]

	return Args{
		ImageURL:        ImageURL,
		ServiceCopyName: args[0],
		BranchName:      BranchName,
		NewServiceName:  NewServiceName,
	}
}

func main() {
	// Logging

	var api = docker.Api{ApiUrl: "http://127.0.0.1:2375"}
	var log = logger.DefaultLogger()
	args := ParseArgs()

	fmt.Printf("%+v\n", args)

	// Try to get a crated service
	service := api.GetService(args.NewServiceName)

	// Check if the service is already created
	if service.Spec.Name == "" {
		log.Warningf("Service %s not created, creating a new one based on develop...", args.BranchName)
		service = api.GetService(args.ServiceCopyName)

		newService := service

		// Change the service informations
		newService.Spec.TaskTemplate.ContainerSpec.Image = args.ImageURL
		newService.Spec.Name = args.NewServiceName
		newService.Spec.Labels["traefik.frontend.rule"] = "Host " + args.BranchName + ".doare.org"

		response := api.CreateService(newService.Spec)
		log.Debugf("Service created with ID: %s", response.ID)
	} else { // Update a created service
		log.Debugf("Updating a existent service with name %s...", args.BranchName)
		newService := service

		// Change the service informations
		newService.Spec.TaskTemplate.ContainerSpec.Image = args.ImageURL
		newService.Spec.Name = args.NewServiceName
		newService.Spec.Labels["traefik.frontend.rule"] = "Host " + args.BranchName + ".doare.org"

		if api.UpdateService(*newService) == true {
			log.Noticef("Service %s updated successfully!", args.BranchName)
		} else {
			log.Errorf("Failure when updating service %s!", args.BranchName)
		}
	}

}
