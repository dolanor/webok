package main

import (
	"fmt"
	"io/ioutil"
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
	c := http.Client{
		Timeout: 30 * time.Second,
	}

	ticker := time.NewTicker(5 * time.Second)
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
	for range ticker.C {
		start := time.Now()
		resp, err := c.Get("https://free.fr")
		end := time.Now()

		if err != nil {
			fmt.Fprintf(w, "%s\t❌\t%v\n", 0*time.Millisecond, err)
			w.Flush()
			systray.SetIcon(noSignalIconData)
			continue
		}
		defer resp.Body.Close()

		systray.SetIcon(signal100IconData)
		dur := end.Sub(start)

		fmt.Fprintf(w, "%s\t✓\n", dur.Truncate(time.Millisecond*10))
		w.Flush()
	}
}
