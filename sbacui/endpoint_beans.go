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
	if err = content.SetOrigin(0, 0); err != nil {
		return err
	}

	for contextName, beans := range beans.Contexts {
		if _, err := h1.Fprintln(content, contextName); err != nil {
			log.Panicln(err)
		}
		for beanName, bean := range beans.Beans {
			if _, err = fmt.Fprintf(content,
				"\t\t%s\n" +
				"\t\t\t\t%s%s\n" +
				"\t\t\t\t%s%s\n" +
				"\t\t\t\t%s%s\n" +
				"\t\t\t\t%s%s\n",
				h3.Sprint(beanName),
				h5.Sprint("Type: "), bean.Type,
				h5.Sprint("Scope: "), bean.Scope,
				h5.Sprint("Aliases: "), strings.Join(bean.Aliases, ", "),
				h5.Sprint("Dependencies: "), strings.Join(bean.Dependencies, ", "),
			); err != nil {
				return err
			}
		}
	}

	return nil
}