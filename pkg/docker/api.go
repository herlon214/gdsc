package docker

import (
	"encoding/json"
	"strconv"

	"github.com/franela/goreq"
	"github.com/herlon214/gdsc/pkg/logger"
	"github.com/herlon214/gdsc/pkg/types"
)

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
func (api *Api) CreateService(spec types.Spec) *types.CreateServiceResponse {
	body, _ := api.Request("POST", "/services/create", spec)

	var response types.CreateServiceResponse
	json.Unmarshal([]byte(body), &response)

	return &response
}

// UpdateService update a docker service based on the given spec
func (api *Api) UpdateService(service types.Service) bool {
	log := logger.DefaultLogger()
	newVersion := strconv.Itoa(service.Version.Index)

	log.Debugf("New version: %s", newVersion)

	_, res := api.Request("POST", "/services/"+service.Spec.Name+"/update?version="+newVersion, service.Spec)

	return res.StatusCode == 200
}

// GetService return the service information
func (api *Api) GetService(nameOrID string) *types.Service {
	body, _ := api.Request("GET", "/services/"+nameOrID, nil)

	var response types.Service
	json.Unmarshal([]byte(body), &response)

	return &response
}
