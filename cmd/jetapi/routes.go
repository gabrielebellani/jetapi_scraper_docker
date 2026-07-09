package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/macsencasaus/jetapi/internal/sites"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api", app.api)

	return mux
}

func (app *application) api(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	countable := false
	defer func() {
		if !countable {
			return
		}
		latency := time.Since(start)
		app.apiCalls.Add(1)
		app.totalLatency.Add(latency.Nanoseconds())
	}()

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	q, err := app.parseAPIQueries(w, r)
	if err != nil {
		app.badRequest(w)
		return
	}
	countable = true

	var jsonResult []byte

	if q.OnlyJP == q.OnlyFR {
		sr, err := sites.Scrape(q)
		if sr == nil {
			app.notFound(w)
			return
		}

		if err != nil {
			app.logErr(fmt.Errorf("Partial Error: %v", err))
		}

		jsonResult, err = json.Marshal(sr)
	} else if q.OnlyJP {
		jpRes, err := sites.ScrapeJetPhotos(q)
		if err != nil {
			app.notFound(w)
			return
		}

		jsonResult, err = json.Marshal(jpRes)
	} else if q.OnlyFR {
		frRes, err := sites.ScrapeFlightRadar(q)
		if err != nil {
			app.notFound(w)
			return
		}

		jsonResult, err = json.Marshal(frRes)
	}

	if err != nil {
		app.serverError(w, fmt.Errorf("Error encoding json: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult)
}
