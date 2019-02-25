package main

import (
	"encoding/json"
	"fmt"
	"github.com/jroimartin/gocui"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const EndpointEnv = "env"

var emptyEnvResponse = EnvResponse{}

type EnvResponse struct {
	ActiveProfiles []string
	PropertySources []EnvPropertySource
}

type EnvPropertySource struct {
	Name string
	Properties EnvProperties
}

type EnvProperties map[string]EnvPropData

type EnvPropData struct {
	Value interface{}
	Origin string
}

func init() {
	log.Println("init env")
	endpointFuncMap[EndpointEnv] = env
}

func env(g *gocui.Gui, link *ActuatorLink, attr map[string]interface{}) error {
	envResponse, err := getEnvResponse(link, attr)
	if err != nil {
		return err
	}

	content, err := g.View(contentView)
	if err != nil {
		return err
	}
	content.Clear()
	if err = content.SetOrigin(0, 0); err != nil {
		return err
	}

	h1.Fprintln(content, "Active Profiles = " + strings.Join(envResponse.ActiveProfiles, ", "))
	for _, propertySource := range envResponse.PropertySources {
		h3.Fprintf(content, "\t\t%s\n", propertySource.Name)
		for key, value := range propertySource.Properties {
			fmt.Fprintf(content, "\t\t\t\t%s = %v\n", h5.Sprint(key), value.Value)
			if value.Origin != "" {
				fmt.Fprintf(content, "\t\t\t\t\t\t%s\n",value.Origin)
			}
		}
	}

	return nil
}

func getEnvResponse(link *ActuatorLink, attr map[string]interface{}) (EnvResponse, error) {
	resp, err := http.Get(link.Href)
	if err != nil {
		return emptyEnvResponse, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return emptyEnvResponse, err
	}
	var envResponse EnvResponse
	err = json.Unmarshal(body, &envResponse)
	if err != nil {
		return emptyEnvResponse, err
	}
	return envResponse, nil
}