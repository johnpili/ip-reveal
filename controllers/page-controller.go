package controllers

import (
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
}

// IndexHandler ...
func (z *PageController) IndexHandler(w http.ResponseWriter, r *http.Request) {
	ipInfo := z.getIPDetails(r)
	x, err := assembler.AssembleTemplate(z.ViewBox, "base.html", "index.html")
	err = x.Execute(w, map[string]string{
		"Title":     "IP Echo",
		"IP":        ipInfo.IP,
		"UserAgent": ipInfo.UserAgent,
	})

	if err != nil {
		log.Panic(err.Error())
	}
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
	}
	return ipInfo
}
