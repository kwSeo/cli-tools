package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type BeansResponse struct {
	Contexts BeansContext
}

type BeansContext map[string]BeansApplication

type BeansApplication struct {
	Beans BeansMap
	ParentId string
}

type BeansMap map[string]Bean

type Bean struct {
	Aliases      []string
	Scope        string
	Type         string
	Resource     string
	Dependencies []string
}

const EndpointBeans = "beans"

func init() {
	log.Println("init beans")
	endpointFuncMap[EndpointBeans] = beansFunc
}

func beansFunc(g *gocui.Gui, link *ActuatorLink, _ map[string]interface{}) error {
	content, err := g.View(contentView)
	if err != nil {
		return err
	}

	response, err := http.Get(link.Href)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil
	}
	var beans BeansResponse
	err = json.Unmarshal(body, &beans)
	if err !=nil {
		return err
	}

	content.Clear()

	contextFont := color.New(color.FgRed, color.Bold)
	for contextName, beans := range beans.Contexts {
		if _, err := contextFont.Fprintln(content, contextName); err != nil {
			log.Panicln(err)
		}
		for beanName, bean := range beans.Beans {
			if _, err = fmt.Fprintf(content,
				"\t\t%s\n" +
				"\t\t\t\t%s%s\n" +
				"\t\t\t\t%s%s\n" +
				"\t\t\t\t%s%s\n" +
				"\t\t\t\t%s%s\n",
				color.CyanString(beanName),
				color.MagentaString("Type: "), bean.Type,
				color.MagentaString("Scope: "), bean.Scope,
				color.MagentaString("Aliases: "), strings.Join(bean.Aliases, ", "),
				color.MagentaString("Dependencies: "), strings.Join(bean.Dependencies, ", "),
			); err != nil {
				log.Panicln(err)
			}
		}
	}

	return nil
}