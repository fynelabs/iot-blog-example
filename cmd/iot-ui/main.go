package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"sort"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/grandcat/zeroconf"
	"github.com/ssimunic/gosensors"
	"golang.org/x/sync/errgroup"
)

func findSensorWithZeroConf() (*zeroconf.ServiceEntry, error) {
	var r *zeroconf.ServiceEntry

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, err
	}

	eg, ctx := errgroup.WithContext(context.Background())

	entries := make(chan *zeroconf.ServiceEntry)
	zcctx, cancel := context.WithTimeout(ctx, time.Second*5)
	eg.Go(func() error {
		for entry := range entries {
			r = entry
			cancel()
		}

		return nil
	})

	eg.Go(func() error {
		err := resolver.Browse(zcctx, "_go_iot_sensor._tcp", "local.", entries)
		if err != nil {
			return err
		}

		<-zcctx.Done()
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return r, nil
}

func getSensorsValueFromZeroConf(sensor *zeroconf.ServiceEntry) (*gosensors.Sensors, error) {
	var client = &http.Client{Timeout: 10 * time.Second}

	var url string
	if len(sensor.AddrIPv4) != 0 {
		url = fmt.Sprintf("http://%v:%v/api", sensor.AddrIPv4[0].String(), sensor.Port)
	} else {
		url = fmt.Sprintf("http://%v:%v/api", sensor.AddrIPv6[0].String(), sensor.Port)
	}
	r, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var sensors gosensors.Sensors
	err = json.NewDecoder(r.Body).Decode(&sensors)
	if err != nil {
		return nil, err
	}

	return &sensors, nil
}

func stringsFromSensors(sensors *gosensors.Sensors) []string {
	var sensorsValue []string
	for chipName, chip := range sensors.Chips {
		for valueName, value := range chip {
			sensorsValue = append(sensorsValue, fmt.Sprintf("%v/%v: %v", chipName, valueName, value))
		}
	}
	sort.Strings(sensorsValue)
	return sensorsValue
}

func main() {
	jsonPtr := flag.Bool("json", false, "get the json content of a remote sensor.")
	flag.Parse()

	remoteSensors, err := findSensorWithZeroConf()
	if err != nil {
		panic(err)
	}

	if *jsonPtr {
		sensors, err := getSensorsValueFromZeroConf(remoteSensors)
		if err != nil {
			panic(err)
		}
		fmt.Println(sensors.JSON())
		return
	}

	a := app.New()
	w := a.NewWindow("IoT go sensors!")
	w.Resize(fyne.NewSize(200, 600))

	hello := widget.NewLabel("List of sensors found on " + remoteSensors.Instance + ":")

	sensorsValue := []string{}

	list := widget.NewList(
		func() int {
			return len(sensorsValue)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(sensorsValue[id])
		},
	)

	go func() {
		for {
			time.Sleep(1 * time.Second)
			sensors, err := getSensorsValueFromZeroConf(remoteSensors)
			if err == nil {
				sensorsValue = stringsFromSensors(sensors)
				list.Refresh()
			}
		}
	}()

	w.SetContent(container.NewBorder(hello, nil, nil, nil, list))

	if _, ok := a.(desktop.App); ok {
		selfManage(a, w)
	}

	w.ShowAndRun()
}
