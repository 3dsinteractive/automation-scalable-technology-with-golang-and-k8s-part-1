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
	// 1. We will check dependency if it is OK, in this case the dependency is Redis
	isAlive := ms.isCacherAlive()
	if !isAlive {
		return false, "Cacher healthcheck failed"
	}

	// 2. If we have other dependency, we will add them here such as
	// isAlive = ms.isMariaDBAlive()
	// if !isAlive {
	// 	return false, "MariaDB healthcheck failed"
	// }

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
		// If Microservice isAlive return !ok, it is because some dependency is failed
		ok, reason := ms.isAlive()
		if !ok {
			// If !ok we will response status 500 error
			ms.responseProbeFailed(c.Response(), reason)
			return nil
		}
		// If ok we will response 200 OK
		ms.responseProbeOK(c.Response())
		return nil
	})
}
