package main

import (
	"encoding/json"
	"fmt"
	"github.com/jroimartin/gocui"
	"io/ioutil"
	"log"
	"net/http"
)

const EndpointMetrics = "metrics"

type MetricData struct {
	Names []string
}

func init() {
	log.Println("init metrics")
	endpointFuncMap[EndpointMetrics] = metrics
}

func metrics(g *gocui.Gui, link *ActuatorLink, _ map[string]interface{}) error {
	v, err := pushSideMenu(g)
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

	var metricData MetricData
	err = json.Unmarshal(body, &metricData)
	if err != nil {
		return err
	}

	metricStr := ""
	for _, name := range metricData.Names {
		metricStr += name + "\n"
	}

	v.Clear()
	v.SetOrigin(0, 0)
	fmt.Fprintln(v, metricStr)

	return nil
}