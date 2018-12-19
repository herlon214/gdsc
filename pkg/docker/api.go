package docker

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/herlon214/gdsc/pkg/http"
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

type Secret struct {
	SecretID   string
	SecretName string
	File       struct {
		Name string
		UID  string
		GID  string
		Mode int
	}
}

type ContainerSpec struct {
	Image     string
	Isolation string
	Env       []string
	Configs   []Config
	Labels    map[string]string
	Mounts    []Mount
	Hosts     []string
	Secrets   []Secret
}

type Placement struct {
	Constraints []string
	Preferences []struct {
		Spread struct {
			SpreadDescriptor string
		}
	}
}

type TaskTemplate struct {
	ContainerSpec ContainerSpec
	ForceUpdate   int
	Runtime       string
	Placement     Placement
}

type Mode struct {
	Replicated struct {
		Replicas int
	}
}

type Port struct {
	Protocol      string
	TargetPort    int
	PublishedPort int
	PublishMode   string
}

type EndpointSpec struct {
	Mode string
	// Ports []Port
	// Not copy the published ports because Traefik will connect to specified container port
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

// Docker service struct
type Service struct {
	Spec    Spec
	Version struct {
		Index int
	}
}

// CreateServiceResponse format of docker api response when create a service
type CreateServiceResponse struct {
	Message  string
	ID       string
	Warnings string
}

type Api struct {
	ApiUrl string
}

// CreateService create a docker service based on the given spec
func (api *Api) CreateService(spec Spec, headers map[string]string) *CreateServiceResponse {
	body, _ := http.Post(api.ApiUrl+"/services/create", spec, headers)

	var response CreateServiceResponse
	json.Unmarshal([]byte(body), &response)

	return &response
}

// CreateRawService create a service from a given json
func (api *Api) CreateRawService(spec Spec) *CreateServiceResponse {
	body, _ := http.Post(api.ApiUrl+"/services/create", spec, nil)

	var response CreateServiceResponse
	json.Unmarshal([]byte(body), &response)

	return &response
}

// UpdateService update a docker service based on the given spec
func (api *Api) UpdateService(service Service, headers map[string]string) bool {
	newVersion := strconv.Itoa(service.Version.Index)

	_, res := http.Post(api.ApiUrl+"/services/"+service.Spec.Name+"/update?version="+newVersion, service.Spec, headers)

	return res.StatusCode == 200
}

// GetService return the service information
func (api *Api) GetService(nameOrID string) *Service {
	body, _ := http.Get(api.ApiUrl + "/services/" + nameOrID)

	var response Service
	json.Unmarshal([]byte(body), &response)

	return &response
}

// GetRawService return the exactly output from docker api
func (api *Api) GetRawService(nameOrID string) string {
	body, _ := http.Get(api.ApiUrl + "/services/" + nameOrID)

	return body
}

// UpdateWithDaemon call os.exec with docker daemon
// because docker service update --with-registry doesn't work yet
func (api *Api) UpdateWithDaemon(service Service) bool {
	err := SystemExec([][]string{
		{
			"docker", "service", "update", "--with-registry-auth", "--image", service.Spec.TaskTemplate.ContainerSpec.Image, service.Spec.Name,
		},
	})

	if err != nil {
		panic(fmt.Errorf("Command failed with: %+v", err))
	}

	return true
}

// SystemExec execute a command in the system
func SystemExec(commands [][]string) error {
	for _, c := range commands {
		cmd := exec.Command(c[0], c[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
