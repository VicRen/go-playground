package main

import (
	"flag"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr = flag.String("addr", ":10997", "The address to listen on for HTTP requests.")

	testingData = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gauge_with_location",
			Help: "Gauge data with location",
		},
		[]string{"latitude", "longitude", "name"})

	testingData2 = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gauge_with_country_code",
			Help: "Gauge data with country code",
		},
		[]string{"cc"})
	locationMap = map[int]*location{
		1: {"40.190632", "116.412144", "beijing"},
		2: {"31.2322758", "121.4692071", "shanghai"},
		3: {"30.2322758", "120.4692071", "ningbo"},
		4: {"29.65262", "91.13775", "lasa"},
		5: {"34.23053", "108.93425", "xian"},
	}
)

func init() {
	prometheus.MustRegister(testingData)
	prometheus.MustRegister(testingData2)
}

func main() {
	flag.Parse()

	go func() {
		for {
			for _, v := range locationMap {
				testingData.WithLabelValues(v.lat, v.lon, v.name).Set(rand.Float64() * 100)
			}
			time.Sleep(5 * time.Second)
		}
	}()
	go func() {
		for {
			testingData2.WithLabelValues("US").Set(rand.Float64() * 100)
			testingData2.WithLabelValues("AE").Set(rand.Float64() * 100)
			testingData2.WithLabelValues("UK").Set(rand.Float64() * 100)
			testingData2.WithLabelValues("JP").Set(rand.Float64() * 100)
			testingData2.WithLabelValues("KR").Set(rand.Float64() * 100)
			testingData2.WithLabelValues("GE").Set(rand.Float64() * 100)
			testingData2.WithLabelValues("FR").Set(rand.Float64() * 100)
			testingData2.WithLabelValues("IR").Set(rand.Float64() * 100)
			testingData2.WithLabelValues("BR").Set(rand.Float64() * 100)
			time.Sleep(10 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))
	if err := http.ListenAndServe(*addr, nil); err != nil {
		panic(err)
	}
}

type location struct {
	lat  string
	lon  string
	name string
}
