package main

import "github.com/labstack/echo"

// HTTPContext implement IContext it is context for HTTP
type HTTPContext struct {
	ms *Microservice
	c  echo.Context
}

// NewHTTPContext is the constructor function for HTTPContext
func NewHTTPContext(ms *Microservice, c echo.Context) *HTTPContext {
	return &HTTPContext{
		ms: ms,
		c:  c,
	}
}
