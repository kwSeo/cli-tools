package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	baseUrlPtr = flag.String(
	"baseUrl",
	"http://localhost:8080/actuator",
	"base URL of Spring Boot Actuator(default: http://localhost:8080/actuator)",
	)
	logFilePtr = flag.String(
		"logFile",
		"",
		"path of log file for debugging",
	)
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
	baseUrl := *baseUrlPtr
	logFile := *logFilePtr
	fmt.Println("BaseURL =", baseUrl)
	fmt.Println("logPath =", logFile)

	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_CREATE | os.O_APPEND | os.O_RDWR, os.ModePerm)
		if err != nil {
			log.Panicln(err)
		}
		defer file.Close()
		log.SetOutput(file)
		log.Println("init logFile")
	}

	resp, err := http.Get(baseUrl)
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

