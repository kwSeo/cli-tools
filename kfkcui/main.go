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
	defaultTopic = "test"
	viewProducer  = "viewProducer"
	viewConsumer  = "viewConsumer"
)

var (
	addrs = flag.String("addr", defaultAddrs, "kafka broker address")
	topic = flag.String("topic", defaultTopic, "kafka topic")
)
type KfkGui struct {
	producer sarama.SyncProducer
	consumer sarama.Consumer
	topic string
}

func NewKfkGui(addrs []string, topic string) *KfkGui {
	producer, err := sarama.NewSyncProducer(addrs, nil)
	utils.MustNotErr(err)

	consumer, err := sarama.NewConsumer(addrs, nil)
	utils.MustNotErr(err)

	return &KfkGui{
		producer: producer,
		consumer: consumer,
		topic: topic,
	}
}

func (k *KfkGui) Close() {
	k.producer.Close()
	k.consumer.Close()
}

func (k *KfkGui) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(viewProducer, 0, 0, maxX/2, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView(viewProducer); err != nil {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Title = "Producer"
		v.Editable = true
	}

	if v, err := g.SetView(viewConsumer, maxX/2, 0, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Consumer"
		v.Autoscroll = true

		go k.consume(g, v)
	}

	return nil
}

func (k *KfkGui) produce(g *gocui.Gui, v *gocui.View) error {
	v.EditNewLine()
	buf := v.ViewBufferLines()
	data := buf[len(buf)-1]
	msg := sarama.ProducerMessage{Topic: k.topic, Value: sarama.StringEncoder(data)}
	if _, _, err := k.producer.SendMessage(&msg); err != nil {
		return err
	}

	return nil
}

func (k *KfkGui) consume(g *gocui.Gui, out *gocui.View) error {
	partition, err := k.consumer.ConsumePartition(k.topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer partition.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	ConsumerLoop:for {
		select {
		case msg := <- partition.Messages():
			g.Update(func(_ *gocui.Gui) error {
				fmt.Fprintf(out, "%s\n", msg.Value)
				return nil
			})
		case <- signals:
			break ConsumerLoop
		}
	}

	return nil
}

func doNothing(_ *gocui.Gui, _ *gocui.View) error {
	return nil
}

func (k *KfkGui) keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding(viewProducer, gocui.KeyEnter, gocui.ModNone, k.produce); err != nil {
		return err
	}
	if err := g.SetKeybinding(viewProducer, gocui.KeyArrowUp, gocui.ModNone, doNothing); err != nil {
		return err
	}
	if err := g.SetKeybinding(viewProducer, gocui.KeyArrowDown, gocui.ModNone, doNothing); err != nil {
		return err
	}
	return nil
}

func quit(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	flag.Parse()
	log.Printf("%s %s\n", *addrs, *topic)

	parsedAddrs := strings.Split(*addrs, addrSeparator)
	kfkGui := NewKfkGui(parsedAddrs, *topic)
	defer kfkGui.Close()

	g, err := gocui.NewGui(gocui.OutputNormal)
	utils.MustNotErr(err)
	defer g.Close()

	g.Cursor = true

	g.SetManager(kfkGui)
	err = kfkGui.keybindings(g)
	utils.MustNotErr(err)

	if err := g.MainLoop(); err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
