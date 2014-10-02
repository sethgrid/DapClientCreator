Sample Client:


package main

import (
	"log"

	"github.com/sendgrid/DapClientCreator/client/user"
)

func main() {
	log.Println("starting")
	client, err := user.New("localhost:9000")
	if err != nil {
		log.Fatalf("Cannot create client: %v", err)
	}

	getter := client.UserGET()

	getter.SetLimit(2) // example, does not make sense in this case
	getter.SetId(180)
	_, err = getter.Do()
	if err != nil {
		log.Println("unable to get user 180 ", err)
	}

	data := getter.Success()

	log.Printf("response: %+v", data)

	// we do not want to collect the response from the *.Do()
	// we can, and we can inspect it, but we drain the response body at that point
	// and it will not be available to the *.Success() call.

	// body, _ := ioutil.ReadAll(resp.Body) // err: is_new_newsletter string vs int
	// log.Printf("raw response %s", body)


}
