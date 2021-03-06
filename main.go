package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	consume "github.com/log-shiper/consume"
	"github.com/log-shiper/g"
	"github.com/log-shiper/httpserver"
	"github.com/log-shiper/produce"
	"github.com/log-shiper/tool"
)

var (
	h bool
	f string
	t string
	b string
	a string
	p string
)

type LogProcess struct {
	read  Reader
	write Writer
	ch    chan string
}

type Reader interface {
	Read(ch chan string)
}

type Writer interface {
	Write(ch chan string)
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: log-shiper  [-f filename] [-t topic] [-b brokers] [-a ip] [-p port] -[h]`)
	flag.PrintDefaults()
	os.Exit(-1)
}

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&a, "a", "127.0.0.1", "log agent ip address")
	flag.StringVar(&b, "b", "", "kafka broker address")
	flag.StringVar(&f, "f", "", "log file name")
	flag.StringVar(&p, "p", "", "log agent port")
	flag.StringVar(&t, "t", "", "topic name")
	flag.Usage = usage
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
	}
	arg := map[string]string{
		"broker": b,
		"file":   f,
		"ip":     a,
		"topic":  t,
		"port":   p,
	}
	n, ret := tool.Argument(arg)
	if !ret {
		logMsg := fmt.Sprintf("%s is null", n)
		log.Print(logMsg)
		flag.Usage()
	}
	r := &produce.ReadFromFile{
		Path: f,
	}
	w := &consume.WriteToKafka{
		Brokers: b,
		Topic:   t,
		MsgKey: g.MsgKey{
			Addr: a,
			Port: p,
		},
	}
	c := make(chan string, 1000)
	lp := &LogProcess{
		read:  r,
		write: w,
		ch:    c,
	}
	go lp.read.Read(lp.ch)
	for i := 1; i <= 5; i++ {
		go lp.write.Write(lp.ch)
	}
	m := httpserver.Monitor{
		StartTime: time.Now(),
		Data:      g.SystemInfo{},
	}
	m.Start(lp.ch, p)
}
