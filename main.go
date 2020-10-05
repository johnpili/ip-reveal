package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/go-zoo/bone"
	"github.com/johnpili/ip-echo/controllers"
	"github.com/johnpili/ip-echo/models"
	"gopkg.in/yaml.v2"
)

// Configurations / Settings
var (
	configuration models.Config
	BuildVersion  = ""
)

func main() {

	pid := os.Getpid()
	err := ioutil.WriteFile("application.pid", []byte(strconv.Itoa(pid)), 0666) // Used to kill this program
	if err != nil {
		log.Print(err)
	}

	var configLocation string
	flag.StringVar(&configLocation, "config", "config.yml", "Set the location of configuration file")
	flag.Parse()

	printProductBanner()

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
	router.HandleFunc("/ip", pageController.TextHandler)

	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
	}

	if configuration.HTTP.IsTLS {
		log.Printf("Server running at https://localhost:%s/\n", port)
		log.Fatal(httpServer.ListenAndServeTLS(configuration.HTTP.ServerCert, configuration.HTTP.ServerKey))
		return
	}
	log.Printf("Server running at http://localhost:%s/\n", port)
	log.Fatal(httpServer.ListenAndServe())
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

func printProductBanner() {
	fmt.Println("")
	fmt.Println("Product: ip.johnpili.com")
	fmt.Println("Version: 2021.1.0")
	fmt.Println("Operating System: ", runtime.GOOS)
	fmt.Println("Architecture: ", runtime.GOARCH)
	fmt.Println("Build Version: ", BuildVersion)
	fmt.Println("System Date/Time: ", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("")
	fmt.Println("")
}
