package main

import (
	"encoding/json"
	"github.com/johnpili/ip-reveal/models"
	"github.com/johnpili/ip-reveal/page"
	"gopkg.in/yaml.v2"
	"html/template"
	"log"
	"net/http"
	"os"
)

func renderPage(w http.ResponseWriter, r *http.Request, vm interface{}, basePath string, filenames ...string) {
	p := vm.(*page.Page)

	if p.Data == nil {
		p.SetData(make(map[string]interface{}))
	}

	if p.ErrorMessages == nil {
		p.ResetErrors()
	}

	if p.UIMapData == nil {
		p.UIMapData = make(map[string]interface{})
	}
	p.UIMapData["basePath"] = basePath

	templateFS := template.Must(template.New("base").ParseFS(views, filenames...))
	err := templateFS.Execute(w, p)
	if err != nil {
		log.Panic(err.Error())
	}
}

func loadConfiguration(a string, b *models.Config) error {
	f, err := os.Open(a)
	if err != nil {
		return err
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(b)
	if err != nil {
		return err
	}
	return nil
}

func loadCountries(a []byte, b *map[string]string) error {
	err := json.Unmarshal(a, b)
	if err != nil {
		return err
	}
	return nil
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
