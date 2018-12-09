package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var baseUrlPtr = flag.String(
	"baseUrl",
	"http://localhost:8080/actuator",
	"base URL of Spring Boot Actuator(default: http://localhost:8080/actuator)",
)

type ActuatorLink struct {
	Href string `json:"href"`
	Templated bool `json:"templated"`
}

type Actuator struct {
	Links map[string]ActuatorLink `json:"_links"`
}

func main() {
	flag.Parse()
	fmt.Printf("BaseURL=%s\n", *baseUrlPtr)

	resp, err := http.Get(*baseUrlPtr)
	if err != nil {
		log.Panicln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicln(err)
	}

	var act Actuator
	err = json.Unmarshal(body, &act)
	if err != nil {
		log.Panicln(err)
	}

	for key, value := range act.Links {
		log.Printf("%s = %v\n", key, value)
	}

	start(&act)
}

