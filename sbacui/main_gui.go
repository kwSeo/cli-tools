package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/kwseo/cli-tools/sbacui/font"
	"log"
	"sort"
	"strconv"
)

const (
	sideView = "side"
	contentView = "content"
)

var (
	sideMenuStack = []string{sideView}
	endpointFuncMap = make(map[string]EndpointFunc)
)

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

	if v, err := g.SetView(sideView, 0, 2, maxX/100*30, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.SelFgColor = gocui.ColorBlack
		v.SelBgColor = gocui.ColorGreen

		for _, endpoint := range ag.Endpoints {
			fmt.Fprintln(v, endpoint)
		}

		g.SetCurrentView(sideView)
		g.SetViewOnTop(sideView)
	}

	for _, name := range sideMenuStack {
		if _, err := g.SetView(name, 0, 2, maxX/100*30, maxY); err != nil {
			return err
		}
	}

	if v, err := g.SetView("content", maxX/100*30, 2, maxX-1, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = false
		v.Autoscroll = false
		fmt.Fprintln(v, "this is content")
	}

	return nil
}

func (ag *ActuatorGui) selectEndpoint(g *gocui.Gui, v *gocui.View) error {
	lines := v.ViewBufferLines()
	_, oy := v.Origin()
	_, cy := v.Cursor()
	y := cy + oy
	if y >= len(lines) {
		return nil
	}

	selected := lines[y]
	endpoint := ag.Actuator.Links[selected]

	if endpointFunc, exist := endpointFuncMap[selected]; exist {
		if err := endpointFunc(g, &endpoint, map[string]interface{}{}); err != nil {
			return err
		}
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


func cursorUp(g *gocui.Gui, v *gocui.View) error {
	ox, oy := v.Origin()
	x, y := v.Cursor()
	if err := v.SetCursor(x, y-1); err != nil && oy > 0 {
		if err := v.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	x, y := v.Cursor()
	if err := v.SetCursor(x, y+1); err != nil {
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	return nil
}

func cursorLeft(g *gocui.Gui, v *gocui.View) error {
	_, err := g.SetCurrentView(contentView)
	return err
}

func cursorRightOnContentView(g *gocui.Gui, v *gocui.View) error {
	_, err := g.SetCurrentView(sideView)
	return err
}

func cursorDownOnContentView(g *gocui.Gui, v *gocui.View) error {
	x, y := v.Origin()
	if err := v.SetOrigin(x, y + 1); err != nil {
		return err
	}
	return nil
}

func cursorUpOnContentView(g *gocui.Gui, v *gocui.View) error {
	x, y := v.Origin()
	if y == 0 {
		return nil
	}
	if err := v.SetOrigin(x, y - 1); err != nil {
		return err
	}
	return nil
}

func (ag *ActuatorGui) keybindings(g *gocui.Gui) error {
	keybind(g, "", gocui.KeyCtrlC, gocui.ModNone, quit)
	keybind(g, sideView, gocui.KeyEnter, gocui.ModNone, ag.selectEndpoint)
	keybind(g, sideView, gocui.KeyArrowDown, gocui.ModNone, cursorDown)
	keybind(g, sideView, gocui.KeyArrowUp, gocui.ModNone, cursorUp)
	keybind(g, sideView, gocui.KeyArrowRight, gocui.ModNone, cursorLeft)
	keybind(g, contentView, gocui.KeyArrowLeft, gocui.ModNone, cursorRightOnContentView)
	keybind(g, contentView, gocui.KeyArrowDown, gocui.ModNone, cursorDownOnContentView)
	keybind(g, contentView, gocui.KeyArrowUp, gocui.ModNone, cursorUpOnContentView)
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func keybind(g *gocui.Gui, name string, key interface{}, mod gocui.Modifier, handler func(*gocui.Gui, *gocui.View) error) {
	if err := g.SetKeybinding(name, key, mod, handler); err != nil {
		log.Panicln(err)
	}
}
func backSideView(g *gocui.Gui, v *gocui.View) error {
	return popSideMenu(g)
}

func pushSideMenu(g *gocui.Gui) (*gocui.View, error) {
	size := len(sideMenuStack)
	name := sideView + strconv.Itoa(size)
	sideMenuStack = append(sideMenuStack, name)
	maxX, maxY := g.Size()
	v, err := g.SetView(name, 0, 2, maxX/100*30, maxY)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}

		v.Highlight = true
		v.SelFgColor = gocui.ColorBlack
		v.SelBgColor = gocui.ColorGreen

		keybind(g, name, gocui.KeyArrowDown, gocui.ModNone, cursorDown)
		keybind(g, name, gocui.KeyArrowUp, gocui.ModNone, cursorUp)
		keybind(g, name, gocui.KeyBackspace2, gocui.ModNone, backSideView)

		if _, err := g.SetCurrentView(name); err != nil {
			return nil, err
		}
		if _, err := g.SetViewOnTop(name); err != nil {
			return nil, err
		}
	}
	return v, nil
}

func popSideMenu(g *gocui.Gui) error {
	top := len(sideMenuStack) - 1
	name := sideMenuStack[top]
	g.DeleteKeybindings(name)
	err := g.DeleteView(name)
	if err != nil {
		return err
	}
	sideMenuStack = sideMenuStack[:top]
	currentTopView := sideMenuStack[top-1]
	if _, err = g.SetCurrentView(currentTopView); err != nil {
		return err
	}
	if _, err := g.SetViewOnTop(currentTopView); err != nil {
		return err
	}
	return nil
}

func currentSideMenu(g *gocui.Gui) (*gocui.View, error) {
	return g.View(sideMenuStack[len(sideMenuStack)-1])
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