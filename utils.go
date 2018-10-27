package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type JSON interface{}

// Get a json decoder for a given requests body
func requestBodyDecoder(request *http.Request) (*json.Decoder, error) {
	// Read body to buffer
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	// Because go lang is a pain in the ass if you read the body then any susequent calls
	// are unable to read the body again....
	request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(body))), nil
}


// Parse the requests body
func parseRequestBody(request *http.Request) (JSON, error) {
	decoder, err := requestBodyDecoder(request)

	if err != nil{
		return nil, err
	}
	var requestPayload JSON
	err = decoder.Decode(&requestPayload)

	if err != nil {
		return nil, err
	}

	return requestPayload, nil
}


func parseRequestBodyAsJsonArray(request *http.Request) ([]JSON, error){
	decoder, err := requestBodyDecoder(request)

	if err != nil{
		return nil, err
	}

	var requestPayload []JSON
	err = decoder.Decode(&requestPayload)

	if err != nil {
		return nil, err
	}

	return requestPayload, nil
}