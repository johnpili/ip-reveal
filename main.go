package main

import (
	"embed"
	_ "embed"
	"flag"
	"github.com/johnpili/ip-reveal/models"
	"github.com/johnpili/ip-reveal/page"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// region Configurations / Settings
var (
	configuration models.Config

	//go:embed views/*
	views embed.FS

	//go:embed countries.json
	countriesJson []byte

	countries map[string]string
)

//endregion

func main() {
	pid := os.Getpid()
	err := os.WriteFile("application.pid", []byte(strconv.Itoa(pid)), 0666) // Used to kill this program
	if err != nil {
		log.Print(err)
	}

	var configLocation string
	flag.StringVar(&configLocation, "config", "config.yml", "Set the location of configuration file")
	flag.Parse()

	err = loadConfiguration(configLocation, &configuration)
	if err != nil {
		log.Fatal(err)
	}

	err = loadCountries(countriesJson, &countries)
	if err != nil {
		log.Fatal(err)
	}

	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/", indexHandler)
	router.HandlerFunc(http.MethodGet, "/json", jsonHandler)
	router.HandlerFunc(http.MethodGet, "/text", textHandler)
	router.HandlerFunc(http.MethodGet, "/txt", textHandler)
	router.HandlerFunc(http.MethodGet, "/ip", textHandler)

	port := strconv.Itoa(configuration.HTTP.Port)
	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
	}

	if configuration.HTTP.IsTLS {
		log.Printf("Server running at https://localhost:%s%s/\n", port, configuration.HTTP.BasePath)
		log.Fatal(httpServer.ListenAndServeTLS(configuration.HTTP.ServerCert, configuration.HTTP.ServerKey))
		return
	}
	log.Printf("Server running at http://localhost:%s%s/\n", port, configuration.HTTP.BasePath)
	log.Fatal(httpServer.ListenAndServe())
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	ipInfo := getIPDetails(r)

	country, ok := countries[ipInfo.IPCountry]
	if !ok {
		country = "Unknown"
	}

	p := page.New()
	p.Title = "IP Reveal"
	p.SetData(map[string]string{
		"IP":        ipInfo.IP,
		"UserAgent": ipInfo.UserAgent,
		"IPCountry": ipInfo.IPCountry,
		"Country":   country,
	})
	renderPage(w, r, p, configuration.HTTP.BasePath, "views/base.html", "views/index.html")
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	ipInfo := getIPDetails(r)
	respondWithJSON(w, ipInfo)
}

func textHandler(w http.ResponseWriter, r *http.Request) {
	ipInfo := getIPDetails(r)
	respondWithPlainText(w, []byte(ipInfo.IP))
}

func getIPDetails(r *http.Request) models.IPInfo {
	ip := ""
	if len(configuration.Extraction.HeaderKey) > 0 {
		ip = r.Header.Get(configuration.Extraction.HeaderKey) // Extract IP from header because we are using reverse proxy example X-Real-Ip
	}

	if len(ip) == 0 { // Fallback
		ip = extractIPAddress(r.RemoteAddr)
	}

	ipInfo := models.IPInfo{
		IP:        ip,
		UserAgent: r.Header.Get("User-Agent"),
		IPCountry: r.Header.Get("CF-IPCountry"),
	}

	if configuration.Extraction.DebugHeader {
		log.Print(r.Header)
	}
	return ipInfo
}

func extractIPAddress(ip string) string {
	if len(ip) > 0 {
		for i := len(ip); i >= 0; i-- {
			offset := len(ip)
			if (i + 1) <= len(ip) {
				offset = i + 1
			}
			if ip[i:offset] == ":" {
				return ip[:i]
			}
		}
	}
	return ip
}
