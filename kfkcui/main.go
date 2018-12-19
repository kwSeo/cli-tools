package main

import (
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/jroimartin/gocui"
	"github.com/kwseo/cli-tools/utils"
	"log"
	"os"
	"os/signal"
	"strings"
)

const (
	addrSeparator = ","
	defaultAddrs  = "127.0.0.1:9092"
	viewLeftSide  = "leftSide"
)

var addrs = flag.String("addr", "", "kafka broker address")

func produce(producer sarama.SyncProducer) {
	log.Println("Start producing")
	msg := sarama.ProducerMessage{Topic: "test", Value: sarama.StringEncoder("SUCCESS")}
	if _, _, err := producer.SendMessage(&msg); err != nil {
		log.Panicln(err)
	}
}

func consume(consumer sarama.Consumer) {
	log.Println("Start consuming")

	partition, err := consumer.ConsumePartition("test", 0, sarama.OffsetNewest)
	if err != nil {
		log.Panicln(err)
	}
	defer partition.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	ConsumerLoop:for {
		select {
		case msg := <- partition.Messages():
			log.Printf("%s\n", msg.Value)
		case <- signals:
			break ConsumerLoop
		}
	}
}

func quit(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}

func layout(g *gocui.Gui) error {
	_, maxY := g.Size()
	if v, err := g.SetView(viewLeftSide, 0, 0, 10, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView(viewLeftSide); err != nil {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, viewLeftSide)
		fmt.Fprintln(v, "wait...")
	}
	return nil
}

func main() {
	flag.Parse()

	if *addrs == "" {
		*addrs = defaultAddrs
	}

	parsedAddrs := strings.Split(*addrs, addrSeparator)
	producer, err := sarama.NewSyncProducer(parsedAddrs, nil)
	utils.MustNotErr(err)
	defer producer.Close()

	consumer, err := sarama.NewConsumer(parsedAddrs, nil)
	utils.MustNotErr(err)
	defer consumer.Close()

	g, err := gocui.NewGui(gocui.OutputNormal)
	utils.MustNotErr(err)
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(layout)ß

	utils.MustNotErr(keybindings(g))

	if err := g.MainLoop(); err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
