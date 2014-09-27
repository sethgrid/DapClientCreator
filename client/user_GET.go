package DapClient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (c *Client) UserGET() *UserGETstruct {
	return &UserGETstruct{httpClient: c.HttpClient, dapAddr: c.DapAddr}
}

type UserGETstruct struct {
	ArgLimit  *int    `json:"limit,omitempty"`
	ArgName   *string `json:"name,omitempty"`
	ArgOffset *int    `json:"offset,omitempty"`
	ArgEmail  *string `json:"email,omitempty"`
	ArgId     *int    `json:"id,omitempty"`

	httpClient *http.Client
	response   *http.Response
	dapAddr    string
}

type UserGETresponse struct {
	Offset string `json:"offset"`
	Email  string `json:"email"`
	Id     string `json:"id"`
	Limit  string `json:"limit"`
	Name   string `json:"name"`
}

func (x *UserGETstruct) Method() string {
	return GET
}

func (x *UserGETstruct) Required() []string {
	return []string{}
}

func (x *UserGETstruct) Location() string {
	return "localhost:9000/api/v1/crud/user"
}

func (x *UserGETstruct) Do() (*http.Response, error) {
	json, err := json.Marshal(x)
	if err != nil {
		// TODO: proper error handling
		log.Fatalf("error marshalling %v", err)
	}
	body := bytes.NewReader(json)

	// location may have parameters in it (/blah/:foo/blah/:bar)
	// these must match to an Arg value on the struct and be replaced.
	l := x.Location()
	strconv.ParseBool("true")
	if x.ArgEmail != nil {
		l = strings.Replace(l, ":email", *x.ArgEmail, -1)
	}
	if x.ArgId != nil {
		l = strings.Replace(l, ":id", strconv.Itoa(*x.ArgId), -1)
	}
	if x.ArgLimit != nil {
		l = strings.Replace(l, ":limit", strconv.Itoa(*x.ArgLimit), -1)
	}
	if x.ArgName != nil {
		l = strings.Replace(l, ":name", *x.ArgName, -1)
	}
	if x.ArgOffset != nil {
		l = strings.Replace(l, ":offset", strconv.Itoa(*x.ArgOffset), -1)
	}

	request, err := http.NewRequest(x.Method(), "http://"+l, body)
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

func (x *UserGETstruct) Success() []UserGETresponse {
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
	response := make([]UserGETresponse, 0)
	err = json.Unmarshal(data, &response)
	if err != nil {
		// TODO: proper error handling
		log.Fatalf("error reading sucesss body %v", err)
	}
	return response
}

func (x *UserGETstruct) Failure() *ErrorResponse {
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

func (x *UserGETstruct) SetEmail(email string) {
	x.ArgEmail = &email
}

func (x *UserGETstruct) SetId(id int) {
	x.ArgId = &id
}

func (x *UserGETstruct) SetLimit(limit int) {
	x.ArgLimit = &limit
}

func (x *UserGETstruct) SetName(name string) {
	x.ArgName = &name
}

func (x *UserGETstruct) SetOffset(offset int) {
	x.ArgOffset = &offset
}
