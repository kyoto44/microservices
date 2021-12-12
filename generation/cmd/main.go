package main

import (
	"encoding/json"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

const (
	numberOfWorkers    = 4
	persistencePostURI = "http://persistence:8001/post"
)

type RandomData struct {
	ID          string       `json:"id"`
	Fingerprint *Fingerprint `json:"fingerprint"`
}

type Fingerprint struct {
	UserAgent string `json:"useragent"`
	IP        string `json:"ip"`
	Country   string `json:"country"`
}

func GenerateData() *RandomData {
	return &RandomData{
		ID: gofakeit.UUID(),
		Fingerprint: &Fingerprint{
			UserAgent: gofakeit.UserAgent(),
			IP:        gofakeit.IPv4Address(),
			Country:   gofakeit.CountryAbr(),
		},
	}
}

func main() {

	for i := 0; i < numberOfWorkers; i++ {
		go Send(time.Duration(1 + rand.Intn(5)))
	}

	log.Info("Generation started successfully")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}

func Send(periodicity time.Duration) {
	for {
		randomDataJSON, err := json.Marshal(GenerateData())
		if err != nil {
			log.Errorf("Can't generate random data, reason: %s", err)
		}

		req := fasthttp.AcquireRequest()
		req.SetBody(randomDataJSON)
		req.Header.SetMethod("POST")
		req.Header.SetContentType("application/json")
		req.SetRequestURI(persistencePostURI)
		res := fasthttp.AcquireResponse()
		if err := fasthttp.Do(req, res); err != nil {
			log.Error(err)
		}
		fasthttp.ReleaseRequest(req)

		time.Sleep(time.Second * periodicity)
	}
}
