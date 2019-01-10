package main

import "github.com/jroimartin/gocui"

const EndpointMetrics = "metrics"

func init() {
	endpointFuncMap[EndpointMetrics] = metrics
}

func metrics(g *gocui.Gui, link *ActuatorLink, _ map[string]interface{}) error {
	v, err := g.View("content")
	if err != nil {
		return err
	}

}