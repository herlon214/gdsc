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
	Domain string `help:"root domain to be used in the traefik host, eg: mycompany.org"`
	Auth   string `help:"registry auth token"`
	APIURL string `help:"docker api url, eg: http://127.0.0.1:2375"`
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
	args := ParseArgs()
	var api docker.Api

	// Check if the user set a different api url
	if args.APIURL != "" {
		api = docker.Api{ApiUrl: args.APIURL}
	} else {
		api = docker.Api{ApiUrl: "http://127.0.0.1:2375"}
	}

	// Try to get a crated service
	service := api.GetService(args.Name)

	// Check if the service is already created
	if service.Spec.Name == "" {
		CreateService(args, api)
	} else { // Update a created service
		UpdateService(args, api, *service)
	}

}

// CreateService creates a new docker service
func CreateService(args Args, api docker.Api) {
	log := logger.DefaultLogger()
	log.Warningf("Service was %s not created yet, creating a new one based on %s ...", args.Name, args.From)
	service := api.GetService(args.From)
	headers := map[string]string{}

	// Check for auth
	if args.Auth != "" {
		headers["X-Registry-Auth"] = args.Auth
	}

	// Change the service informations
	service.Spec.TaskTemplate.ContainerSpec.Image = args.Image
	service.Spec.Name = args.Name

	// Check if must set a Traefik rule
	if args.Domain != "" {
		service.Spec.Labels["traefik.frontend.rule"] = "Host: " + args.Name + "." + args.Domain
	}

	response := api.CreateService(service.Spec, headers)
	log.Debugf("Service created with ID: %s", response.ID)
}

// UpdateService updates a service with new docker image
func UpdateService(args Args, api docker.Api, service docker.Service) {
	log := logger.DefaultLogger()
	headers := map[string]string{}
	log.Debugf("Updating a existent service with name %s ...", args.Name)

	// Check for auth
	if args.Auth != "" {
		headers["X-Registry-Auth"] = args.Auth
	}

	// Change the service informations
	service.Spec.TaskTemplate.ContainerSpec.Image = args.Image

	if api.UpdateService(service, headers) == true {
		log.Noticef("Service %s updated successfully!", args.Name)
	} else {
		log.Errorf("Failure when updating service %s!", args.Name)

		// Exit with error status
		os.Exit(1)
	}
}
