package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type WeatherResponse struct {
	CurrentWeather struct {
		Temperature float64 `json:"temperature"`
	} `json:"current_weather"`
}

var cities = map[string]string{
	"Warszawa (Polska)":  "latitude=52.2297&longitude=21.0122",
	"Paryż (Francja)":    "latitude=48.8566&longitude=2.3522",
	"Tokio (Japonia)":    "latitude=35.6895&longitude=139.6917",
}

func main() {

	health := flag.Bool("health", false, "Wykonaj healthcheck")
	flag.Parse()

	if *health {
		resp, err := http.Get("http://localhost:8080")
		if err != nil || resp.StatusCode != 200 {
			os.Exit(1)
		}
		os.Exit(0) 
	}

	author := "KACPER SOKÓŁ" 
	port := "8080"
	startupTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("Data uruchomienia: %s | Autor: %s | Port TCP nasłuchu: %s\n", startupTime, author, port)


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		selectedCity := r.URL.Query().Get("city")
		var weatherInfo string

	
		if coords, exists := cities[selectedCity]; exists {
			apiURL := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?%s&current_weather=true", coords)
			resp, err := http.Get(apiURL)
			if err == nil {
				defer resp.Body.Close()
				body, _ := io.ReadAll(resp.Body)
				var weather WeatherResponse
				json.Unmarshal(body, &weather)
				weatherInfo = fmt.Sprintf("Aktualna temperatura dla %s to: %.1f°C", selectedCity, weather.CurrentWeather.Temperature)
			} else {
				weatherInfo = "Błąd pobierania pogody."
			}
		}

		html := `
		<!DOCTYPE html>
		<html>
		<head><title>Pogoda</title><meta charset="utf-8"></head>
		<body>
			<h2>Sprawdź pogodę</h2>
			<form method="GET">
				<select name="city">
					<option value="" disabled selected>Wybierz miasto...</option>
					<option value="Warszawa (Polska)">Warszawa (Polska)</option>
					<option value="Paryż (Francja)">Paryż (Francja)</option>
					<option value="Tokio (Japonia)">Tokio (Japonia)</option>
				</select>
				<button type="submit">Sprawdź</button>
			</form>
			<br>
			<h3>{{.Weather}}</h3>
		</body>
		</html>`
		
		tmpl, _ := template.New("webpage").Parse(html)
		tmpl.Execute(w, struct{ Weather string }{weatherInfo})
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
