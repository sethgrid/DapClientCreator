package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/sendgrid/dap/meta"
)

var metaUrl string
var templateDir string

func init() {
	flag.StringVar(&metaUrl, "metaUrl", "http://localhost:9000/api/v1/_meta", "Pass in the url of the meta endpoint for Dap")
	flag.StringVar(&templateDir, "templateDir", "client", "The directory relative to parent where to place templated files")
}

func main() {
	flag.Parse()
	log.Println("Getting meta data from " + metaUrl)
	resp, err := http.Get(metaUrl)
	if err != nil {
		log.Fatal("unable to get metaUrl ", err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("unable to read response body ", err)
	}
	metas, err := ParseMeta(data)

	createBaseTemplateFiles()

	for _, m := range metas {
		createTemplateFile(m)
	}

	gofmt()
}

func gofmt() {
	_, err := exec.Command("go", "fmt", "./...").Output()
	if err != nil {
		log.Println("error running go fmt ", err)
	}
}

func ParseMeta(data []byte) ([]meta.Meta, error) {
	metas := make([]meta.Meta, 0)
	err := json.Unmarshal(data, &metas)
	return metas, err
}

func createBaseTemplateFiles() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("unable to get working directory ", err)
	}
	path := dir + "/" + templateDir
	log.Println("removing templateDir")

	err = os.RemoveAll(path)
	if !os.IsNotExist(err) && err != nil {
		log.Fatal("unable to remove templateDir ", err)
	}

	log.Println("creating templateDir")
	err = os.Mkdir(path, 0777)
	if os.IsNotExist(err) && err != nil {
		log.Fatal("unable to create templateDir ", err)
	}

}

func createTemplateFile(m meta.Meta) {
	funcMap := template.FuncMap{
		"title":                protectKeywordsUpper,
		"lower":                protectKeywordsLower,
		"structHelper":         structHelper,
		"structResponseHelper": structResponseHelper,
		"accessorHelper":       accessorHelper,
		"requiredList":         requiredList,
		"locationHelper":       locationHelper,
	}

	tmpl, err := template.New("myTemplate").Funcs(funcMap).Parse(OUTPUT_FILE)
	if err != nil {
		log.Fatal("error creating template ", err)
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("unable to get working directory ", err)
	}
	path := dir + "/" + templateDir
	err = os.Mkdir(path+"/"+m.Title, 0777)
	if err != nil && os.IsExist(err) != true {
		log.Fatal("unable to create client subpackage directory ", err)
	}
	log.Println(fmt.Sprintf("creating %s/%s_%s.go", m.Title, m.Title, m.Method))

	f, err := os.Create(fmt.Sprintf("%s/%s/%s_%s.go", path, m.Title, m.Title, m.Method))
	defer f.Close()
	if err != nil {
		log.Fatal("unable to create template file ", err)
	}

	err = tmpl.Execute(f, m)
	if err != nil {
		log.Fatal("unable to execute template ", err)
	}

	clientFile := path + "/" + m.Title + "/client.go"
	if _, err := os.Stat(clientFile); os.IsNotExist(err) {
		log.Println("creating local client.go")
		f, err := os.OpenFile(clientFile, os.O_WRONLY|os.O_CREATE, 0655)
		defer f.Close()
		if err != nil {
			log.Fatal("unable to create client.go ", err)
		}

		content := client_file(m.Title)
		_, err = f.Write([]byte(content))
		if err != nil {
			log.Fatal("unable write to client.go ", err)
		}
	}
}

func requiredList(m meta.Meta) string {
	tmp := make([]string, len(m.Required))
	for i, v := range m.Required {
		tmp[i] = fmt.Sprintf(`"%s"`, v)
	}
	return strings.Join(tmp, ",")
}

func structHelper(m meta.Meta) string {
	var returnString string

	for name, info := range m.Properties {
		returnString += fmt.Sprintf("Arg%s *%s `json:\"%s,omitempty\"`\n",
			strings.Title(protectKeywords(name)), Mysql2GoType(info.DataType), strings.ToLower(protectKeywords(name)))
	}
	return returnString
}

func structResponseHelper(m meta.Meta) string {
	var returnString string

	for name, info := range m.Properties {
		returnString += fmt.Sprintf("%s %s `json:\"%s\"`\n",
			strings.Title(protectKeywords(name)), Mysql2GoType(info.DataType), strings.ToLower(protectKeywords(name)))
	}
	return returnString
}

func accessorHelper(m meta.Meta) string {
	var returnString string

	for name, info := range m.Properties {
		returnString += fmt.Sprintf(`
			func (x *%s%sstruct) Set%s(%s %s){
				x.Arg%s = &%s
			}
			`,
			protectKeywordsLower(m.Title),
			strings.Title(m.Method),

			protectKeywordsUpper(name),
			protectKeywords(name),
			Mysql2GoType(info.DataType),

			protectKeywordsUpper(name),
			protectKeywords(name))
	}
	return returnString
}

func protectKeywordsUpper(w string) string {
	return strings.Title(protectKeywords(w))
}

func protectKeywordsLower(w string) string {
	return strings.ToLower(protectKeywords(w))
}

func protectKeywords(w string) string {
	switch w {
	// whatever we go with, we need to make sure the structs stay exportable (no prepend _)
	case "type":
		return w + "_"
	case "package":
		return w + "_"
	case "func":
		return w + "_"
	case "var":
		return w + "_"
		// TODO - fix my keywords in dap to not use Offset and Limit. Conflicts started as noted below.
		// case "offset":
		// 	return strings.Title(w) + "_"
		// case "limit":
		// 	return strings.Title(w) + "_"
	}
	return w
}

func locationHelper(m meta.Meta) string {
	// many of the possible returnStrings require strconv, but not all.
	// make sure that strconv is used to avoid import errors
	returnString := "strconv.ParseBool(\"true\")"

	// above wont work due to "value" needing to be dynamic each run.
	// will have to do multiple replacement attempts
	for k, v := range m.Properties {
		returnString += `
		if x.Arg` + protectKeywordsUpper(k) + ` != nil{`
		switch v.DataType {
		case "timestamp":
			returnString += `
			l = strings.Replace(l, ":` + k + `", strconv.FormatFloat(float64(*x.Arg` + protectKeywordsUpper(k) + `), 'f', -1, 32), -1)`
		case "varchar":
			returnString += `
			l = strings.Replace(l, ":` + k + `", *x.Arg` + protectKeywordsUpper(k) + `, -1)`
		case "tinyint":
			// canot formatbool due to json decoder issue (int=>bool)
			returnString += `
			l = strings.Replace(l, ":` + k + `", strconv.Itoa(*x.Arg` + protectKeywordsUpper(k) + `), -1)`
		case "smallint":
			// canot formatbool due to json decoder issue (int=>bool)
			returnString += `
			l = strings.Replace(l, ":` + k + `", strconv.Itoa(*x.Arg` + protectKeywordsUpper(k) + `), -1)`
		case "int":
			returnString += `
			l = strings.Replace(l, ":` + k + `", strconv.Itoa(*x.Arg` + protectKeywordsUpper(k) + `), -1)`
		default:
			// TODO: add notifications/logging to discover this early
			returnString += `
			// unknown datatype ` + v.DataType + `
			//l = strings.Replace(l, ":` + k + `", *x.Arg` + protectKeywordsUpper(k) + `, -1)`
		}
		returnString += `
		}`
	}

	return returnString
}

func Mysql2GoType(MysqlType string) string {
	var GoType string

	switch MysqlType {
	case "timestamp":
		GoType = "float32"
	case "varchar":
		GoType = "string"
	case "tinyint":
		GoType = "int" // cannot yet use bool. json decoding issue type int=>bool... custom decoder?
	case "smallint":
		GoType = "int" // cannot yet use bool. json decoding issue type int=>bool... custom decoder?
	case "int":
		GoType = "int"
	default:
		GoType = "string"
	}

	return GoType
}

/*
{
        "description": "Custom Endpoint POST Example",
        "location": "api/v1/custom/endpoint",
        "method": "POST",
        "notes": "Used as an example endpoint. Will accept data.",
        "primary": "id",
        "properties": {
            "sample_property": {
                "description": "example property",
                "type": "varchar"
            }
        },
        "required": [],
        "responseSchema": {
            "InsertedID": "int",
            "Message": "string"
        },
        "title": "",
        "type": "mysql"
    }
*/

const OUTPUT_FILE = `package {{ lower .Title}}
import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"strconv"
)


type {{lower .Title}}{{title .Method}}struct struct{
	{{ structHelper . }}

	httpClient *http.Client
	response   *http.Response
	dapAddr    string
}

func (c *Client) {{title .Title}}{{title .Method}}() *{{lower .Title}}{{title .Method}}struct {
	// avoid missing import error
	if false{
		strings.Title("foo")
	}
	return &{{lower .Title}}{{title .Method}}struct{httpClient: c.HttpClient, dapAddr: c.DapAddr}
}

type {{title .Title}}{{title .Method}}response struct{
	{{ structResponseHelper . }}
}

func (x *{{lower .Title}}{{title .Method}}struct) Method() string{
	return {{ .Method }}
}

func (x *{{lower .Title}}{{title .Method}}struct) Required() []string{
	return []string{ {{ requiredList . }} }
}

func (x *{{lower .Title}}{{title .Method}}struct) Location() string {
	return "{{ .Location }}"
}

func (x *{{lower .Title}}{{title .Method}}struct) Do() (*http.Response, error) {
	json, err := json.Marshal(x)
	if err != nil {
		// TODO: proper error handling
		log.Fatalf("error marshalling %v", err)
	}
	body := bytes.NewReader(json)

	// location may have parameters in it (/blah/:foo/blah/:bar)
	// these must match to an Arg value on the struct and be replaced.
	l := x.Location()
	{{ locationHelper . }}

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

func (x *{{lower .Title}}{{title .Method}}struct) Success() []{{title .Title}}{{title .Method}}response {
	if x.response == nil {
		return nil
	}

	// get the response body and put it back (as reading drains the response)
	data, err := ioutil.ReadAll(x.response.Body)
	if err != nil {
		// TODO: proper error handling
		log.Fatalf("error reading body %v", err)
	}

	response := make([]{{title .Title}}{{title .Method}}response, 0)

	x.response.Body = ioutil.NopCloser(bytes.NewReader(data))

	dataReader := bytes.NewReader(data)
	// decoder is better than unmarshall at handling streams
	decoder := json.NewDecoder(dataReader)
    err = decoder.Decode(&response)

	if err != nil {
		// TODO: proper error handling
		log.Fatalf("error reading sucesss body %v", err)
	}
	return response
}

func (x *{{lower .Title}}{{title .Method}}struct) Failure() *ErrorResponse {
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
{{ accessorHelper . }}

`

func client_file(packageName string) string {
	return `package ` + protectKeywords(packageName) + `

import (
	"net/http"
)

type Client struct {
	DapAddr    string
	HttpClient *http.Client
}

func New(dapUrl string) (*Client, error) {
	client := &Client{DapAddr: dapUrl, HttpClient: &http.Client{}}
	return client, nil
}

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type ErrorResponse struct {
	Error string
}
`
}
