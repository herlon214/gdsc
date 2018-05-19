package docker

import (
	"encoding/json"
	"strconv"

	"github.com/herlon214/gdsc/pkg/http"
	"github.com/herlon214/gdsc/pkg/logger"
)

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

type Config struct {
	ConfigID   string
	ConfigName string
	File       struct {
		Name string
		UID  string
		GID  string
		Mode int
	}
}

type Mount struct {
	Type   string
	Source string
	Target string
}

type ContainerSpec struct {
	Image     string
	Isolation string
	Env       []string
	Configs   []Config
	Labels    map[string]string
	Mounts    []Mount
}

type TaskTemplate struct {
	ContainerSpec ContainerSpec
	ForceUpdate   int
	Runtime       string
}

type Mode struct {
	Replicated struct {
		Replicas int
	}
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
	EndpointSpec   struct {
		Mode string
	}
}

// Docker service struct
type Service struct {
	Spec    Spec
	Version struct {
		Index int
	}
}

// CreateServiceResponse format of docker api response when create a service
type CreateServiceResponse struct {
	message string
	ID      string
}

type Api struct {
	ApiUrl string
}

// CreateService create a docker service based on the given spec
func (api *Api) CreateService(spec Spec) *CreateServiceResponse {
	body, _ := http.Post(api.ApiUrl+"/services/create", spec)

	var response CreateServiceResponse
	json.Unmarshal([]byte(body), &response)

	return &response
}

// UpdateService update a docker service based on the given spec
func (api *Api) UpdateService(service Service) bool {
	log := logger.DefaultLogger()
	newVersion := strconv.Itoa(service.Version.Index)

	log.Debugf("New version: %s", newVersion)

	_, res := http.Post(api.ApiUrl+"/services/"+service.Spec.Name+"/update?version="+newVersion, service.Spec)

	return res.StatusCode == 200
}

// GetService return the service information
func (api *Api) GetService(nameOrID string) *Service {
	body, _ := http.Get(api.ApiUrl + "/services/" + nameOrID)

	var response Service
	json.Unmarshal([]byte(body), &response)

	return &response
}
