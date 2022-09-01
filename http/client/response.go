package client

import (
	"io"
	"bytes"
	"net/http"
	"encoding/json"
	"github.com/kelchy/go-lib/log"
)

// Res - response instance from http call
type Res struct {
	Response	http.Response
	Error		error
	log		log.Log
}

// HTML - method to return the html content of response
func (r Res) HTML() (string, error) {
	var data bytes.Buffer
	if r.Error != nil {
		return data.String(), r.Error
	}
	defer r.Response.Body.Close()
	_, e := io.Copy(&data, r.Response.Body)
	return data.String(), e
}

// JSON - method to return the json content of response
func (r Res) JSON() (interface{}, error) {
	var data interface{}
	if r.Error != nil {
		return data, r.Error
	}
	defer r.Response.Body.Close()
	e := json.NewDecoder(r.Response.Body).Decode(&data)
	return data, e
}
