package DapClient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func (c *Client) SettingsGET() *SettingsGETstruct {
	return &SettingsGETstruct{httpClient: c.HttpClient, dapAddr: c.DapAddr}
}

type SettingsGETstruct struct {
	ArgUser_id *int    `json:"user_id,omitempty"`
	ArgEnabled *bool   `json:"enabled,omitempty"`
	ArgId      *int    `json:"id,omitempty"`
	ArgLimit   *int    `json:"limit,omitempty"`
	ArgOffset  *int    `json:"offset,omitempty"`
	ArgSetting *string `json:"setting,omitempty"`

	httpClient *http.Client
	response   *http.Response
	dapAddr    string
}

type SettingsGETresponse struct {
	User_id string `json:"user_id"`
	Enabled string `json:"enabled"`
	Id      string `json:"id"`
	Limit   string `json:"limit"`
	Offset  string `json:"offset"`
	Setting string `json:"setting"`
}

func (x *SettingsGETstruct) Method() string {
	return GET
}

func (x *SettingsGETstruct) Required() []string {
	return []string{}
}

func (x *SettingsGETstruct) Location() string {
	return "localhost:9000/api/v1/crud/settings"
}

func (x *SettingsGETstruct) Do() (*http.Response, error) {
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

func (x *SettingsGETstruct) Success() []SettingsGETresponse {
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
	response := make([]SettingsGETresponse, 0)
	err = json.Unmarshal(data, &response)
	if err != nil {
		// TODO: proper error handling
		log.Fatalf("error reading sucesss body %v", err)
	}
	return response
}

func (x *SettingsGETstruct) Failure() *ErrorResponse {
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

func (x *SettingsGETstruct) SetLimit(limit int) {
	x.ArgLimit = &limit
}

func (x *SettingsGETstruct) SetOffset(offset int) {
	x.ArgOffset = &offset
}

func (x *SettingsGETstruct) SetSetting(setting string) {
	x.ArgSetting = &setting
}

func (x *SettingsGETstruct) SetUser_id(user_id int) {
	x.ArgUser_id = &user_id
}

func (x *SettingsGETstruct) SetEnabled(enabled bool) {
	x.ArgEnabled = &enabled
}

func (x *SettingsGETstruct) SetId(id int) {
	x.ArgId = &id
}
