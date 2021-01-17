// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/parnurzeal/gorequest"
)

// IRequester is interface to connect to HTTP endpoint
type IRequester interface {
	Get(path string, params map[string]string) (string, error)
	Post(path string, params map[string]string) (string, error)
	PostJSON(path string, body interface{}) (string, error)
	Put(path string, params map[string]string) (string, error)
	PutJSON(path string, body interface{}) (string, error)
	Delete(path string, params map[string]string) (string, error)
}

// Requester implement IRequester
type Requester struct {
	ms      *Microservice
	baseURL string
	req     *gorequest.SuperAgent
	timeout time.Duration
}

// NewRequester return new Requester
func NewRequester(baseURL string, timeout time.Duration, ms *Microservice) *Requester {
	return &Requester{
		baseURL: baseURL,
		ms:      ms,
		timeout: timeout,
	}
}

func (rqt *Requester) cloneR() *gorequest.SuperAgent {
	r := rqt.req
	if r == nil {
		r = gorequest.New()
		rqt.req = r
	}
	// Timeout is relative to time.Now so we need to set every time
	r.Timeout(rqt.timeout)
	return r.Clone()
}

// Get request using HTTP GET
func (rqt *Requester) Get(path string, params map[string]string) (string, error) {

	url := fmt.Sprint(rqt.baseURL, path)

	r := rqt.cloneR()

	r = r.Get(url)
	if params != nil {
		for key, value := range params {
			r = r.Param(key, value)
		}
	}

	res, body, errs := r.End()
	if len(errs) > 0 {
		return "", errs[0]
	}

	if res.StatusCode >= 400 {
		return body, fmt.Errorf(res.Status)
	}

	return body, nil
}

// Delete request using HTTP DELETE
func (rqt *Requester) Delete(path string, params map[string]string) (string, error) {

	url := fmt.Sprint(rqt.baseURL, path)

	r := rqt.cloneR()
	r = r.Delete(url)
	if params != nil {
		for key, value := range params {
			r = r.Param(key, value)
		}
	}

	res, body, errs := r.End()
	if len(errs) > 0 {
		return "", errs[0]
	}

	if res.StatusCode >= 400 {
		return body, fmt.Errorf(res.Status)
	}

	return body, nil
}

// Post request using HTTP POST
func (rqt *Requester) Post(path string, params map[string]string) (string, error) {

	u := fmt.Sprint(rqt.baseURL, path)

	r := rqt.cloneR()

	r = r.Post(u)

	if params != nil {
		postData := url.Values{}
		for key, value := range params {
			postData.Add(key, value)
		}
		postDataStr := postData.Encode()

		r = r.Send(postDataStr)
	}

	res, body, errs := r.End()
	if len(errs) > 0 {
		return "", errs[0]
	}

	if res.StatusCode >= 400 {
		return body, fmt.Errorf(res.Status)
	}

	return body, nil
}

// PostJSON request using HTTP POST with JSON body
func (rqt *Requester) PostJSON(path string, jsonBody interface{}) (string, error) {

	url := fmt.Sprint(rqt.baseURL, path)

	r := rqt.cloneR()

	r = r.Post(url)

	if jsonBody != nil {
		r = r.Send(jsonBody)
	}

	res, body, errs := r.End()
	if len(errs) > 0 {
		return "", errs[0]
	}

	if res.StatusCode >= 400 {
		return body, fmt.Errorf(res.Status)
	}

	return body, nil
}

// Put request using HTTP PUT
func (rqt *Requester) Put(path string, params map[string]string) (string, error) {

	u := fmt.Sprint(rqt.baseURL, path)

	r := rqt.cloneR()

	r = r.Put(u)

	if params != nil {
		postData := url.Values{}
		for key, value := range params {
			postData.Add(key, value)
		}
		postDataStr := postData.Encode()
		r = r.Send(postDataStr)
	}

	res, body, errs := r.End()
	if len(errs) > 0 {
		return "", errs[0]
	}

	if res.StatusCode >= 400 {
		return body, fmt.Errorf(res.Status)
	}

	return body, nil
}

// PutJSON request using HTTP PUT with JSON body
func (rqt *Requester) PutJSON(path string, jsonBody interface{}) (string, error) {

	url := fmt.Sprint(rqt.baseURL, path)

	r := rqt.cloneR()

	r = r.Put(url)

	if jsonBody != nil {
		r = r.Send(jsonBody)
	}

	res, body, errs := r.End()
	if len(errs) > 0 {
		return "", errs[0]
	}

	if res.StatusCode >= 400 {
		return body, fmt.Errorf(res.Status)
	}

	return body, nil
}
