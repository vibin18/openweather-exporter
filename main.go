package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"
	"github.com/vibin18/openweather-exporter/gather"
	"log"
	"net/http"
	"os"
)

type Opts struct {
	ServerBind   string `long:"bind"             env:"SERVER_BIND"       description:"Server address."     default:":8080"`
	CityName     string `long:"city"             env:"CITY_NAME"         description:"City name from openweather list."     default:"Böblingen, DE"`
	ApiKey       string `long:"apikey"           env:"API_KEY"           description:"Key for openweather API." required:"true" `
	ScheduleTime string `long:"cron"             env:"SCHEDULES"         description:"e.g. '*/2 * * * *'  '@midnight' '@every 1h30m'. "     default:"@every 1m"`
}

var (
	Args Opts
)

func startHttpServer() {

	http.Handle("/metrics", promhttp.Handler())
	log.Println("starting metric server on ",Args.ServerBind)
	log.Fatal(http.ListenAndServe(Args.ServerBind, nil))
}

func initArgParser() {
	argparser := flags.NewParser(&Args, flags.Default)
	_, err := argparser.Parse()
	// check if there is a parse error
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println()
			argparser.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}
}

func main() {
	initArgParser()
	log.Println("starting openweather-exporter")
	var new gather.Sig
	c := cron.New(
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger),
			cron.Recover(cron.DefaultLogger),
		))
	var _, err = c.AddFunc(Args.ScheduleTime, func() {
		log.Println("starting to Gather weather data from open weather api.")
		new.Fetchdata(Args.CityName, Args.ApiKey)
		log.Println("gathering completed.")

	})
	if err != nil {
		log.Fatal(err)
	}

	c.Start()

	if err := prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Subsystem: "weather",
			Name:      "openweather_temperature",
			Help:      "current weather temprature",
		},
		func() float64 { return float64(new.GetWeatherTemp()) },
	)); err == nil {
		log.Println("openweather_temperature registered")
	}

	if err := prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Subsystem: "weather",
			Name:      "openweather_pressure",
			Help:      "current weather pressure",
		},
		func() float64 { return float64(new.GetWeatherPressure()) },
	)); err == nil {
		log.Println("openweather_pressure registered")
	}

	if err := prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Subsystem: "weather",
			Name:      "openweather_humidity",
			Help:      "current weather humidity",
		},
		func() float64 { return float64(new.GetWeatherHumidity()) },
	)); err == nil {
		log.Println("openweather_humidity registered")
	}

	if err := prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Subsystem: "weather",
			Name:      "openweather_feelslike",
			Help:      "current weather feelslikey",
		},
		func() float64 { return float64(new.GetWeatherFeelsLike()) },
	)); err == nil {
		log.Println("openweather_feelslike registered")
	}

	startHttpServer()
}
