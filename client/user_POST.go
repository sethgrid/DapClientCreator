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

func (c *Client) UserPOST() *UserPOSTstruct {
	return &UserPOSTstruct{httpClient: c.HttpClient, dapAddr: c.DapAddr}
}

type UserPOSTstruct struct {
	ArgEmail *string `json:"email,omitempty"`
	ArgId    *int    `json:"id,omitempty"`
	ArgName  *string `json:"name,omitempty"`

	httpClient *http.Client
	response   *http.Response
	dapAddr    string
}

type UserPOSTresponse struct {
	Email string `json:"email"`
	Id    string `json:"id"`
	Name  string `json:"name"`
}

func (x *UserPOSTstruct) Method() string {
	return POST
}

func (x *UserPOSTstruct) Required() []string {
	return []string{}
}

func (x *UserPOSTstruct) Location() string {
	return "localhost:9000/api/v1/crud/user"
}

func (x *UserPOSTstruct) Do() (*http.Response, error) {
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
	if x.ArgId != nil {
		l = strings.Replace(l, ":id", strconv.Itoa(*x.ArgId), -1)
	}
	if x.ArgName != nil {
		l = strings.Replace(l, ":name", *x.ArgName, -1)
	}
	if x.ArgEmail != nil {
		l = strings.Replace(l, ":email", *x.ArgEmail, -1)
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

func (x *UserPOSTstruct) Success() []UserPOSTresponse {
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
	response := make([]UserPOSTresponse, 0)
	err = json.Unmarshal(data, &response)
	if err != nil {
		// TODO: proper error handling
		log.Fatalf("error reading sucesss body %v", err)
	}
	return response
}

func (x *UserPOSTstruct) Failure() *ErrorResponse {
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

func (x *UserPOSTstruct) SetEmail(email string) {
	x.ArgEmail = &email
}

func (x *UserPOSTstruct) SetId(id int) {
	x.ArgId = &id
}

func (x *UserPOSTstruct) SetName(name string) {
	x.ArgName = &name
}
