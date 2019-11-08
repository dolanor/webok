package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"github.com/getlantern/systray"
	"github.com/markbates/pkger"
)

var (
	signal100IconData []byte
	noSignalIconData  []byte
)

func loadIcons() {
	s100f, err := pkger.Open("/asset/img/nm-signal-100.png")
	if err != nil {
		panic(err)
	}
	nosf, err := pkger.Open("/asset/img/nm-no-connection.png")
	if err != nil {
		panic(err)
	}
	defer s100f.Close()

	noSignalIconData, err = ioutil.ReadAll(nosf)
	if err != nil {
		panic(err)
	}

	signal100IconData, err = ioutil.ReadAll(s100f)
	if err != nil {
		panic(err)
	}
}

func onReady() {
	systray.SetIcon(noSignalIconData)
}

func main() {
	loadIcons()
	go systray.Run(onReady, nil)

	ticker := time.NewTicker(5 * time.Second)
	w := &tabwriter.Writer{}
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	for range ticker.C {
		start := time.Now()
		resp, err := http.Get("https://free.fr")
		end := time.Now()

		if err != nil {
			log.Println(err)
			systray.SetIcon(noSignalIconData)
			continue
		}
		defer resp.Body.Close()

		systray.SetIcon(signal100IconData)
		dur := end.Sub(start)

		fmt.Fprintf(w, "%s\tall good :)\n", dur.Truncate(time.Millisecond*10))
		w.Flush()
	}
}
