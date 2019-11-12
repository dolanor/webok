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

func loadIcons() (noSignalIconData, signal100IconData []byte, err error) {
	s100f, err := pkger.Open("/asset/img/nm-signal-100.png")
	if err != nil {
		return nil, nil, err
	}
	nosf, err := pkger.Open("/asset/img/nm-no-connection.png")
	if err != nil {
		return nil, nil, err
	}
	defer s100f.Close()

	noSignalIconData, err = ioutil.ReadAll(nosf)
	if err != nil {
		return nil, nil, err
	}

	signal100IconData, err = ioutil.ReadAll(s100f)
	if err != nil {
		return nil, nil, err
	}

	return noSignalIconData, signal100IconData, nil
}

func onReady(noSignalIconData []byte, quitCh chan struct{}) func() {
	return func() {
		systray.SetIcon(noSignalIconData)
		quit := systray.AddMenuItem("Quit", "Quit")
		<-quit.ClickedCh
		log.Println("quitting")
		quitCh <- struct{}{}
	}
}

func main() {
	noSignalIconData, signal100IconData, err := loadIcons()
	if err != nil {
		panic(err)
	}

	quitCh := make(chan struct{})
	go systray.Run(onReady(noSignalIconData, quitCh), nil)
	c := http.Client{
		Timeout: 5 * time.Second,
	}

	ticker := time.NewTicker(30 * time.Second)
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
	for range ticker.C {
		select {
		case <-quitCh:
			goto out
		default:
		}

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
out:
	ticker.Stop()
}
