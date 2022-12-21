package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"

	"github.com/grandcat/zeroconf"
	"github.com/ssimunic/gosensors"
)

//go:embed index.html
var indexPage []byte

func serveIndexPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(indexPage)
}

func main() {
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		sensors, err := gosensors.NewFromSystem()
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, sensors.JSON())
	})
	http.HandleFunc("/", serveIndexPageHandler)

	server, err := zeroconf.Register("GoIoTSensor", "_go_iot_sensor._tcp", "local.", 8080, nil, nil)
	if err != nil {
		panic(err)
	}
	defer server.Shutdown()

	log.Println("Server running on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
