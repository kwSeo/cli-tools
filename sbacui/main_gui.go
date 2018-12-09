package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/kwseo/cli-tools/sbacui/font"
	"log"
	"sort"
)

var endpointFuncMap = make(map[string]EndpointFunc)

type EndpointFunc func(*gocui.Gui, *ActuatorLink, map[string]interface{}) error

type ActuatorGui struct {
	Endpoints []string
	Actuator *Actuator
}

func NewActuatorGui(actuator *Actuator) *ActuatorGui {
	var endpoints []string
	for endpoint := range actuator.Links {
		endpoints = append(endpoints, endpoint)
	}

	sort.Strings(endpoints)

	return &ActuatorGui{
		Endpoints: endpoints,
		Actuator: actuator,
	}
}

func (ag *ActuatorGui) Layout(g *gocui.Gui) error {
	title := font.Red(" Spring Boot Actuator CLI")
	maxX, maxY := g.Size()

	if v, err := g.SetView("title", 0, 0, maxX-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, title)
	}

	if v, err := g.SetView("side", 0, 2, maxX/100*30, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelFgColor = gocui.ColorBlack
		v.SelBgColor = gocui.ColorGreen

		for _, endpoint := range ag.Endpoints {
			fmt.Fprintln(v, endpoint)
		}
	}

	if v, err := g.SetView("sideDetail", 0, 2, maxX/100*30, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelFgColor = gocui.ColorBlack
		v.SelBgColor = gocui.ColorGreen
	}

	if v, err := g.SetView("content", maxX/100*30, 2, maxX-1, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Autoscroll = true
		fmt.Fprintln(v, "this is content")
	}

	g.SetCurrentView("side")
	g.SetViewOnTop("side")
	return nil
}

func (ag *ActuatorGui) selectEndpoint(g *gocui.Gui, v *gocui.View) error {
	lines := v.ViewBufferLines()
	_, y := v.Cursor()
	if y >= len(lines) {
		return nil
	}

	selected := lines[y]
	endpoint := ag.Actuator.Links[selected]

	if endpointFunc, exist := endpointFuncMap[selected]; exist {
		endpointFunc(g, &endpoint, map[string]interface{}{})
	} else {
		v, err := g.View("content")
		if err != nil {
			return err
		}
		v.Clear()
		fmt.Fprintln(v, "Not Implemented")
	}

	return nil
}


func (ag *ActuatorGui) cursorUp(g *gocui.Gui, v *gocui.View) error {
	ox, oy := v.Origin()
	x, y := v.Cursor()
	if err := v.SetCursor(x, y-1); err != nil && oy > 0 {
		if err := v.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}
	return nil
}

func (ag *ActuatorGui) cursorDown(g *gocui.Gui, v *gocui.View) error {
	x, y := v.Cursor()
	if err := v.SetCursor(x, y+1); err != nil {
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	return nil
}

func keybind(g *gocui.Gui, name string, key interface{}, mod gocui.Modifier, handler func(*gocui.Gui, *gocui.View) error) {
	if err := g.SetKeybinding(name, key, mod, handler); err != nil {
		log.Panicln(err)
	}
}

func (ag *ActuatorGui) keybindings(g *gocui.Gui) error {
	keybind(g, "", gocui.KeyCtrlC, gocui.ModNone, quit)
	keybind(g, "side", gocui.KeyArrowDown, gocui.ModNone, ag.cursorDown)
	keybind(g, "side", gocui.KeyArrowUp, gocui.ModNone, ag.cursorUp)
	keybind(g, "side", gocui.KeyEnter, gocui.ModNone, ag.selectEndpoint)
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func start(act *Actuator) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	ag := NewActuatorGui(act)

	g.SetManager(ag)

	if err := ag.keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}