package DapClient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func (c *Client) CustomFooPOST() *CustomFooPOSTstruct {
	return &CustomFooPOSTstruct{httpClient: c.HttpClient, dapAddr: c.DapAddr}
}

type CustomFooPOSTstruct struct {
	ArgSample_property *string `json:"sample_property,omitempty"`

	httpClient *http.Client
	response   *http.Response
	dapAddr    string
}

type CustomFooPOSTresponse struct {
	Sample_property string `json:"sample_property"`
}

func (x *CustomFooPOSTstruct) Method() string {
	return POST
}

func (x *CustomFooPOSTstruct) Required() []string {
	return []string{}
}

func (x *CustomFooPOSTstruct) Location() string {
	return "api/v1/custom/endpoint"
}

func (x *CustomFooPOSTstruct) Do() (*http.Response, error) {
	json, err := json.Marshal(x)
	if err != nil {
		// TODO: proper error handling
		log.Fatalf("error marshalling %v", err)
	}
	body := bytes.NewReader(json)
	//request, err := http.NewRequest(x.Method(), x.dapAddr+x.Location(), body)
	request, err := http.NewRequest(x.Method(), "http://"+x.Location(), body)
	if err != nil {
		// TODO: proper error handling
		log.Fatalf("error with new request %v", err)
	}
	response, err := x.httpClient.Do(request)
	if err != nil {
		// TODO: proper error handling
		log.Fatalf("error with response %v", err)
	}
	x.response = response
	return response, nil
}

func (x *CustomFooPOSTstruct) Success() []CustomFooPOSTresponse {
	if x.response == nil {
		return nil
	}

	// get the response body and put it back (as reading drains the response)
	data, err := ioutil.ReadAll(x.response.Body)
	x.response.Body = ioutil.NopCloser(bytes.NewReader(data))

	if err != nil {
		// TODO: proper error handling
		log.Fatalf("error reading body %v", err)
	}
	response := make([]CustomFooPOSTresponse, 0)
	err = json.Unmarshal(data, &response)
	if err != nil {
		// TODO: proper error handling
		log.Fatalf("error reading sucesss body %v", err)
	}
	return response
}

func (x *CustomFooPOSTstruct) Failure() *ErrorResponse {
	if x.response == nil {
		return nil
	}
	// read in the body and put it back
	data, err := ioutil.ReadAll(x.response.Body)
	x.response.Body = ioutil.NopCloser(bytes.NewReader(data))

	if err != nil {
		// TODO: proper error handling
		log.Fatalf("error reading failure body %v", err)
	}
	failure := &ErrorResponse{}
	err = json.Unmarshal(data, failure)
	if err != nil {
		// TODO: proper error handling
		log.Fatalf("error reading sucesss body %v", err)
	}
	return failure
}

// accessor functions

func (x *CustomFooPOSTstruct) SetSample_property(sample_property string) {
	x.ArgSample_property = &sample_property
}
