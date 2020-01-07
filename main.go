package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	rice "github.com/GeertJohan/go.rice"
	"github.com/go-zoo/bone"
	"github.com/johnpili/ip-echo/controllers"
	"github.com/johnpili/ip-echo/models"
	"gopkg.in/yaml.v2"
)

var configuration models.Config

func main() {

	pid := os.Getpid()
	err := ioutil.WriteFile("application.pid", []byte(strconv.Itoa(pid)), 0666) // Used to kill this program
	if err != nil {
		log.Print(err)
	}

	var configLocation string
	flag.StringVar(&configLocation, "config", ".config.yml", "Set the location of configuration file")
	flag.Parse()
	configuration = loadConfiguration(configLocation)

	port := strconv.Itoa(configuration.HTTP.Port)
	if os.Getenv("ASPNETCORE_PORT") != "" {
		port = os.Getenv("ASPNETCORE_PORT") // Override port if deployed in IIS
	}

	viewBox := rice.MustFindBox("views")
	staticBox := rice.MustFindBox("static")

	pageController := controllers.PageController{
		ViewBox:       viewBox,
		Configuration: &configuration,
	}

	staticFileServer := http.StripPrefix("/static/", http.FileServer(staticBox.HTTPBox()))

	router := bone.New()
	router.Handle("/static/", staticFileServer)
	router.HandleFunc("/", pageController.IndexHandler)
	router.HandleFunc("/json", pageController.JSONHandler)
	router.HandleFunc("/text", pageController.TextHandler)
	router.HandleFunc("/txt", pageController.TextHandler)
	log.Fatal(http.ListenAndServe(":"+port, router)) // Start HTTP Server
}

func loadConfiguration(c string) models.Config {
	f, err := os.Open(c)
	if err != nil {
		log.Fatal(err.Error())
	}

	var configuration models.Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatal(err.Error())
	}
	return configuration
}
