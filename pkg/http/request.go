package http

import (
	"encoding/json"

	"github.com/herlon214/gdsc/pkg/logger"
	"github.com/parnurzeal/gorequest"
)

func Get(uri string) (string, gorequest.Response) {
	log := logger.DefaultLogger()
	// log.Debugf("[GET] -> %s", uri)

	request := gorequest.New()
	resp, body, err := request.Get(uri).End()

	if err != nil {
		log.Critical(err)
	}

	return body, resp
}

func Post(uri string, body interface{}) (string, gorequest.Response) {
	log := logger.DefaultLogger()
	// log.Debugf("[POST] -> %s", uri)

	bodyJSON, _ := json.Marshal(body)
	// log.Debugf("Body %+v", string(bodyJSON))

	request := gorequest.New()
	resp, bodyRes, err := request.Post(uri).Send(string(bodyJSON)).End()

	// log.Debug(bodyRes)

	if err != nil {
		log.Critical(err)
	}

	return bodyRes, resp
}
