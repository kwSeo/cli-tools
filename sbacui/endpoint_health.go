package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/tidwall/pretty"
	"io/ioutil"
	"log"
	"net/http"
)

const EndpointHealth = "health"

func init() {
	log.Println("init health")
	endpointFuncMap[EndpointHealth] = health
}

func health(g *gocui.Gui, link *ActuatorLink, _ map[string]interface{}) error {
	contentV, err := g.View("content")
	if err != nil {
		return err
	}

	resp, err := http.Get(link.Href)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	prettyBody := pretty.Pretty(body)

	contentV.Clear()
	if err = contentV.SetOrigin(0, 0); err != nil {
		return err
	}
	fmt.Fprintln(contentV, string(prettyBody))
	return nil
}
