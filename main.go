package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/franela/goreq"
	"github.com/op/go-logging"
)

type ModeReplicated struct {
	Replicas int
}

type Mode struct {
	Replicated ModeReplicated
}

type UpdateConfig struct {
	Parallelism     int
	Delay           int
	FailureAction   string
	Monitor         int
	MaxFailureRatio int
	Order           string
}

type RollbackConfig struct {
	Parallelism     int
	FailureAction   string
	Monitor         int
	MaxFailureRatio int
	Order           string
}

type Network struct {
	Target string
}

type EndpointSpec struct {
	Mode string
}

type ContainerSpec struct {
	Image string
}

type TaskTemplate struct {
	ContainerSpec ContainerSpec
	ForceUpdate   int
	Runtime       string
}

// Docker service spec struct
type Spec struct {
	Name           string
	Labels         map[string]string
	TaskTemplate   TaskTemplate
	Mode           Mode
	UpdateConfig   UpdateConfig
	RollbackConfig RollbackConfig
	Networks       []Network
	EndpointSpec   EndpointSpec
}

type ServiceVersion struct {
	Index int
}

// Docker service struct
type Service struct {
	Spec    Spec
	Version ServiceVersion
}

type DockerApi struct {
	ApiUrl string
}

type ServiceUpdateQueryString struct {
	version int
}

type Args struct {
	ImageURL        string
	BranchName      string
	ServiceCopyName string
	NewServiceName  string
	ServiceVersion  int
}

// Request do a docker api request
func (dockerApi *DockerApi) Request(method string, path string, body interface{}) (string, *goreq.Response) {
	log := GetLogger()

	fullPath := dockerApi.ApiUrl + path

	log.Debugf("[%s] -> %s", method, fullPath)

	bodyJSON, _ := json.Marshal(body)
	log.Debugf("Body %+v", string(bodyJSON))

	res, err := goreq.Request{
		Method: method,
		Uri:    fullPath,
		Body:   body,
	}.Do()

	responseBody, _ := res.Body.ToString()

	log.Debug(responseBody)

	if err != nil {
		log.Critical(err)
	}

	return responseBody, res
}

// CreateServiceResponse format of docker api response when create a service
type CreateServiceResponse struct {
	message string
	ID      string
}

// CreateService create a docker service based on the given spec
func (dockerApi *DockerApi) CreateService(spec Spec) *CreateServiceResponse {
	body, _ := dockerApi.Request("POST", "/services/create", spec)

	var response CreateServiceResponse
	json.Unmarshal([]byte(body), &response)

	return &response
}

// UpdateService update a docker service based on the given spec
func (dockerApi *DockerApi) UpdateService(service Service) bool {
	log := GetLogger()
	newVersion := strconv.Itoa(service.Version.Index)

	log.Debugf("New version: %s", newVersion)

	_, res := dockerApi.Request("POST", "/services/"+service.Spec.Name+"/update?version="+newVersion, service.Spec)

	return res.StatusCode == 200
}

// GetService return the service information
func (dockerApi *DockerApi) GetService(nameOrID string) *Service {
	body, _ := dockerApi.Request("GET", "/services/"+nameOrID, nil)

	var response Service
	json.Unmarshal([]byte(body), &response)

	return &response
}

// GetLogger return a default logger
func GetLogger() *logging.Logger {
	var log = logging.MustGetLogger("service-copy")
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level}%{color:reset} %{message}`,
	)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	logging.SetBackend(backend2Formatter)

	return log
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
	var api = DockerApi{ApiUrl: "http://127.0.0.1:2375"}
	var log = GetLogger()
	args := ParseArgs()

	fmt.Printf("%+v\n", args)

	// Try to get a crated service
	service := api.GetService(args.NewServiceName)

	// Check if the service is already created
	if service.Spec.Name == "" {
		log.Debugf("Service %s not created, creating a new one based on develop...", args.BranchName)
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
