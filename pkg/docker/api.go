package docker

import (
	"encoding/json"
	"strconv"

	"github.com/franela/goreq"
	"github.com/herlon214/gdsc/pkg/logger"
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

type ServiceUpdateQueryString struct {
	version int
}

// CreateServiceResponse format of docker api response when create a service
type CreateServiceResponse struct {
	message string
	ID      string
}

type Api struct {
	ApiUrl string
}

// Request do a docker api request
func (api *Api) Request(method string, path string, body interface{}) (string, *goreq.Response) {
	log := logger.DefaultLogger()

	fullPath := api.ApiUrl + path

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

// CreateService create a docker service based on the given spec
func (api *Api) CreateService(spec Spec) *CreateServiceResponse {
	body, _ := api.Request("POST", "/services/create", spec)

	var response CreateServiceResponse
	json.Unmarshal([]byte(body), &response)

	return &response
}

// UpdateService update a docker service based on the given spec
func (api *Api) UpdateService(service Service) bool {
	log := logger.DefaultLogger()
	newVersion := strconv.Itoa(service.Version.Index)

	log.Debugf("New version: %s", newVersion)

	_, res := api.Request("POST", "/services/"+service.Spec.Name+"/update?version="+newVersion, service.Spec)

	return res.StatusCode == 200
}

// GetService return the service information
func (api *Api) GetService(nameOrID string) *Service {
	body, _ := api.Request("GET", "/services/"+nameOrID, nil)

	var response Service
	json.Unmarshal([]byte(body), &response)

	return &response
}
