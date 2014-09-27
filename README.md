Sample Client:


package main

import (
	"io/ioutil"
	"log"

	"github.com/sendgrid/DapClientCreator/client"
)

func main() {
	log.Println("starting")
	client, err := DapClient.New("localhost:9000")
	if err != nil {
		log.Fatalf("Cannot create client: %v", err)
	}

	log.Println("client created. attempting post...")
	setter := client.CustomFooPOST()
	setter.SetSample_property("this is my sample property")
	setterResponse, err := setter.Do()
	if err != nil {
		log.Println("error with CustomFooPOST ", err)
	}
	log.Println(setterResponse.StatusCode)
	setterResponseBody, _ := ioutil.ReadAll(setterResponse.Body)
	log.Printf("%s\n", setterResponseBody)

	getter := client.CustomFooGET()
	getter.SetSample_property("this is my sample property")
	getterResponse, err := getter.Do()
	if err != nil {
		log.Println("error with CustomFooGET ", err)
	}
	log.Println(getterResponse.StatusCode)
	getterResponseBody, _ := ioutil.ReadAll(getterResponse.Body)
	log.Printf("%s\n", getterResponseBody)

	log.Println("Get Settings...")

	get := client.SettingsGET()

	responseA, err := get.Do()
	if err != nil {
		log.Fatal(err)
	}

	if responseA.StatusCode == 200 {
		res := get.Success()

		for _, successCase := range res {
			log.Printf("Success: User ID (%v) Enabled (%v) Setting (%v)\n", successCase.User_id, successCase.Enabled, successCase.Setting)
		}
	} else {
		res := get.Failure()
		log.Printf("Failure: %s", res.Error)
	}

	log.Println("Read response body...")

	content, err := ioutil.ReadAll(responseA.Body)
	log.Printf("content: %s %v", content, err)

	log.Println("Post Settings...")

	post := client.SettingsPOST()

	post.SetEnabled(true)
	post.SetId(13)
	post.SetSetting("some other setting")
	post.SetUser_id(1)

	responseB, err := post.Do()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Status Response POST: %v", responseB.Status)
	log.Println("Read Response Body...")
	content, err = ioutil.ReadAll(responseB.Body)
	log.Printf("body: %s \nErr: %v", content, err)
}
