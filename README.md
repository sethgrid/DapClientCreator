Sample Client:


package main

import (
	"io/ioutil"
	"log"

	//"github.com/sendgrid/dapclient"
	"github.com/sendgrid/DapClientCreator/client"
)

func main() {
	log.Println("starting")
	client, err := DapClient.New("localhost:9000")
	if err != nil {
		log.Fatalf("Cannot create client: %v", err)
	}

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