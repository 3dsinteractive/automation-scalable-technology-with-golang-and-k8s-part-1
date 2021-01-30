// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func (ms *Microservice) isCacherAlive() bool {
	if ms.cacher == nil {
		return true
	}

	ms.Log("MS", "Perform healthcheck on Cacher")
	err := ms.cacher.Healthcheck()
	if err != nil {
		return false
	}
	return true
}

func (ms *Microservice) isAlive() (bool, string) {
	isAlive := ms.isCacherAlive()
	if !isAlive {
		return false, "Cacher healthcheck failed"
	}

	return true, ""
}

func (ms *Microservice) responseProbeOK(resp *echo.Response) {
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("ok"))
}

func (ms *Microservice) responseProbeFailed(resp *echo.Response, reason string) {
	errMsg := "Healthcheck failed because of " + reason
	resp.WriteHeader(http.StatusInternalServerError)
	resp.Write([]byte(errMsg))
}

// RegisterLivenessProbeEndpoint register endpoint for liveness probe
func (ms *Microservice) RegisterLivenessProbeEndpoint(path string) {
	ms.echo.GET(path, func(c echo.Context) error {
		ok, reason := ms.isAlive()
		if !ok {
			ms.responseProbeFailed(c.Response(), reason)
			return nil
		}
		ms.responseProbeOK(c.Response())
		return nil
	})
}
