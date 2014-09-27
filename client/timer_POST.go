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

func (c *Client) TimerPOST() *TimerPOSTstruct {
	return &TimerPOSTstruct{httpClient: c.HttpClient, dapAddr: c.DapAddr}
}

type TimerPOSTstruct struct {
	ArgIp         *int     `json:"ip,omitempty"`
	ArgPlace_id   *string  `json:"place_id,omitempty"`
	ArgTimer_ms   *int     `json:"timer_ms,omitempty"`
	ArgCreated_at *float32 `json:"created_at,omitempty"`
	ArgId         *int     `json:"id,omitempty"`

	httpClient *http.Client
	response   *http.Response
	dapAddr    string
}

type TimerPOSTresponse struct {
	Place_id   string `json:"place_id"`
	Timer_ms   string `json:"timer_ms"`
	Created_at string `json:"created_at"`
	Id         string `json:"id"`
	Ip         string `json:"ip"`
}

func (x *TimerPOSTstruct) Method() string {
	return POST
}

func (x *TimerPOSTstruct) Required() []string {
	return []string{"place_id", "timer_ms", "created_at"}
}

func (x *TimerPOSTstruct) Location() string {
	return "localhost:9000/api/v1/crud/timer"
}

func (x *TimerPOSTstruct) Do() (*http.Response, error) {
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
	if x.ArgIp != nil {
		l = strings.Replace(l, ":ip", strconv.Itoa(*x.ArgIp), -1)
	}
	if x.ArgPlace_id != nil {
		l = strings.Replace(l, ":place_id", *x.ArgPlace_id, -1)
	}
	if x.ArgTimer_ms != nil {
		l = strings.Replace(l, ":timer_ms", strconv.Itoa(*x.ArgTimer_ms), -1)
	}
	if x.ArgCreated_at != nil {
		l = strings.Replace(l, ":created_at", strconv.FormatFloat(float64(*x.ArgCreated_at), 'f', -1, 32), -1)
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

func (x *TimerPOSTstruct) Success() []TimerPOSTresponse {
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
	response := make([]TimerPOSTresponse, 0)
	err = json.Unmarshal(data, &response)
	if err != nil {
		// TODO: proper error handling
		log.Fatalf("error reading sucesss body %v", err)
	}
	return response
}

func (x *TimerPOSTstruct) Failure() *ErrorResponse {
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

func (x *TimerPOSTstruct) SetCreated_at(created_at float32) {
	x.ArgCreated_at = &created_at
}

func (x *TimerPOSTstruct) SetId(id int) {
	x.ArgId = &id
}

func (x *TimerPOSTstruct) SetIp(ip int) {
	x.ArgIp = &ip
}

func (x *TimerPOSTstruct) SetPlace_id(place_id string) {
	x.ArgPlace_id = &place_id
}

func (x *TimerPOSTstruct) SetTimer_ms(timer_ms int) {
	x.ArgTimer_ms = &timer_ms
}
