package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"io/ioutil"
	"net/http"
)

const EndpointHealth = "health"

func init() {
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

	contentV.Clear()
	fmt.Fprintln(contentV, string(body))
	return nil
}
