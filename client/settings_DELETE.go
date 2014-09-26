package DapClient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func (c *Client) SettingsDELETE() *SettingsDELETEstruct {
	return &SettingsDELETEstruct{httpClient: c.HttpClient, dapAddr: c.DapAddr}
}

type SettingsDELETEstruct struct {
	ArgEnabled *bool   `json:"enabled,omitempty"`
	ArgId      *int    `json:"id,omitempty"`
	ArgLimit   *int    `json:"limit,omitempty"`
	ArgSetting *string `json:"setting,omitempty"`
	ArgUser_id *int    `json:"user_id,omitempty"`

	httpClient *http.Client
	response   *http.Response
	dapAddr    string
}

type SettingsDELETEresponse struct {
	Setting string `json:"setting"`
	User_id string `json:"user_id"`
	Enabled string `json:"enabled"`
	Id      string `json:"id"`
	Limit   string `json:"limit"`
}

func (x *SettingsDELETEstruct) Method() string {
	return DELETE
}

func (x *SettingsDELETEstruct) Required() []string {
	return []string{"id", "limit"}
}

func (x *SettingsDELETEstruct) Location() string {
	return "localhost:9000/api/v1/crud/settings"
}

func (x *SettingsDELETEstruct) Do() (*http.Response, error) {
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

func (x *SettingsDELETEstruct) Success() []SettingsDELETEresponse {
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
	response := make([]SettingsDELETEresponse, 0)
	err = json.Unmarshal(data, &response)
	if err != nil {
		// TODO: proper error handling
		log.Fatalf("error reading sucesss body %v", err)
	}
	return response
}

func (x *SettingsDELETEstruct) Failure() *ErrorResponse {
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

func (x *SettingsDELETEstruct) SetId(id int) {
	x.ArgId = &id
}

func (x *SettingsDELETEstruct) SetLimit(limit int) {
	x.ArgLimit = &limit
}

func (x *SettingsDELETEstruct) SetSetting(setting string) {
	x.ArgSetting = &setting
}

func (x *SettingsDELETEstruct) SetUser_id(user_id int) {
	x.ArgUser_id = &user_id
}

func (x *SettingsDELETEstruct) SetEnabled(enabled bool) {
	x.ArgEnabled = &enabled
}
