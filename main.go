package main

//This script listens on a given TCP port for
//HTTP REST Get messages than scraps the given
//Open vSwtich entry and gives back the stats
//in Prometheus compatible format

//Written by Megyo @ LeanNet

import (
	"log"
	"net/http"

	"github.com/biwwy0/ovs-exporter/ovs"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//the TCP port that this scripts listens
var listenPort string = ":9981"

func handler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	if target == "" {
		http.Error(w, "Bad request!\nCorrect format is: http://"+r.Host+"/metrics?target=<targetIP>", 400)
		return
	}
	c := OvsPromCollector{
		iface:        target,
		port:      ovs.OvsDefaultPort,
		ovsReader: ovs.CliDumpReader,
	}
	registry := prometheus.NewRegistry()
	registry.MustRegister(c)
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/metrics", handler)
	log.Fatal(http.ListenAndServe(listenPort, nil))
}
