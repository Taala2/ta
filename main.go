package main

import (
	"html/template"
	"net/http"
	"strconv"
	"log"
	"fmt"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func convertWeight(value float64, fromUnit, toUnit string) float64 {
	conversion := map[string]float64{
		"miligram": 1,
		"gram":     1000,
		"kilogram": 1000000,
		"ounce":    28349.5,
		"pound":    453592,
	}

	return value * conversion[fromUnit] / conversion[toUnit]

}

func convertLength(value float64, fromUnit, toUnit string) float64 {
	conversion := map[string]float64{
		"milimeter":  1,
		"centimeter": 10,
		"meter":      1000,
		"kilometer":  1000000,
		"inch":       25.4,
		"foot":       304.8,
		"yard":       914.4,
		"mile":       1609344,
	}

	return value * conversion[fromUnit] / conversion[toUnit]
}

func convertTemperature(value float64, fromUnit, toUnit string) float64 {
	if fromUnit == toUnit {
		return value
	}

	switch fromUnit {
	case "celsius":
		if toUnit == "fahrenheit" {
			return value*9/5 + 32
		}else if toUnit == "kelvin" {
			return value + 273.15
		}
	case "fahrenheit":
		if toUnit == "celsius" {
			return (value - 32) * 5 / 9
		}else if toUnit == "kelvin" {
			return (value - 32) * 5 / 9 + 273.15
		}
	case "kelvin":
		if toUnit == "celsius" {
			return value - 273.15
		}else if toUnit == "fahrenheit" {
			return (value - 273.15) * 9 / 5 + 32
		}
	}
	return value
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func lengthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		value, _ := strconv.ParseFloat(r.FormValue("value"), 64)
		fromUnit := r.FormValue("fromUnit")
		toUnit := r.FormValue("toUnit")
		result := convertLength(value, fromUnit, toUnit)
		renderTemplate(w, "length.html", map[string]interface{}{
			"Result":   result,
			"Value":    value,
			"FromUnit": fromUnit,
			"ToUnit":   toUnit,
		})
		return
	}
	renderTemplate(w, "length.html", nil)
}

func weightHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		value, _ := strconv.ParseFloat(r.FormValue("value"), 64)
		fromUnit := r.FormValue("fromUnit")
		toUnit := r.FormValue("toUnit")
		result := convertWeight(value, fromUnit, toUnit)
		renderTemplate(w, "weight.html", map[string]interface{}{
			"Result":   result,
			"Value":    value,
			"FromUnit": fromUnit,
			"ToUnit":   toUnit,
		})
		return
	}
	renderTemplate(w, "weight.html", nil)
}

func temperatureHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		value, _ := strconv.ParseFloat(r.FormValue("value"), 64)
		fromUnit := r.FormValue("fromUnit")
		toUnit := r.FormValue("toUnit")
		result := convertTemperature(value, fromUnit, toUnit)
		renderTemplate(w, "temperature.html", map[string]interface{}{
			"Result":   result,
			"Value":    value,
			"FromUnit": fromUnit,
			"ToUnit":   toUnit,
		})
		return
	}
	renderTemplate(w, "temperature.html", nil)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/length", lengthHandler)
	http.HandleFunc("/weight", weightHandler)
	http.HandleFunc("/temperature", temperatureHandler)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}