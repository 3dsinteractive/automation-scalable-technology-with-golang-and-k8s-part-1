package main

// IContext is the context for service
type IContext interface {
	Log(message string)
	Param(name string) string
	QueryParam(name string) string
	ReadInput() string
	Response(responseCode int, responseData interface{})
}
