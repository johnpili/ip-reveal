package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gorilla/sessions"
	"github.com/johnpili/ip-echo/assembler"
	"github.com/johnpili/ip-echo/models"
)

// PageController ...
type PageController struct {
	Store         *sessions.CookieStore
	ViewBox       *rice.Box
	Configuration *models.Config
	Countries     map[string]string
}

// IndexHandler ...
func (z *PageController) IndexHandler(w http.ResponseWriter, r *http.Request) {
	ipInfo := z.getIPDetails(r)
	x, err := assembler.AssembleTemplate(z.ViewBox, "base.html", "index.html")

	country, ok := z.Countries[ipInfo.IPCountry]
	if !ok {
		country = "Unknown"
	}

	err = x.Execute(w, map[string]string{
		"Title":     "IP Echo",
		"IP":        ipInfo.IP,
		"UserAgent": ipInfo.UserAgent,
		"IPCountry": ipInfo.IPCountry,
		"Country":   country,
	})

	if err != nil {
		log.Panic(err.Error())
	}
}

// JSONHandler ...
func (z *PageController) JSONHandler(w http.ResponseWriter, r *http.Request) {
	ipInfo := z.getIPDetails(r)
	respondWithJSON(w, ipInfo)
}

// TextHandler ...
func (z *PageController) TextHandler(w http.ResponseWriter, r *http.Request) {
	ipInfo := z.getIPDetails(r)
	respondWithPlainText(w, []byte(ipInfo.IP))
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

func (z *PageController) getIPDetails(r *http.Request) models.IPInfo {
	ip := ""
	if len(z.Configuration.Extraction.HeaderKey) > 0 {
		ip = r.Header.Get(z.Configuration.Extraction.HeaderKey) // Extract IP from header because we are using reverse proxy example X-Real-Ip
	}

	if len(ip) == 0 { // Fallback
		ip = extractIPAddress(r.RemoteAddr)
	}

	ipInfo := models.IPInfo{
		IP:        ip,
		UserAgent: r.Header.Get("User-Agent"),
		IPCountry: r.Header.Get("CF-IPCountry"),
	}

	if z.Configuration.Extraction.DebugHeader {
		log.Print(r.Header)
	}

	return ipInfo
}

func respondWithJSON(w http.ResponseWriter, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(200)
	w.Write(response)
}

func respondWithPlainText(w http.ResponseWriter, payload []byte) {
	w.WriteHeader(200)
	w.Write(payload)
}
