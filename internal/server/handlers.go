package server

import (
	"net/http"

	"utm-fwd/internal/utm"
)

type forwardRequest struct {
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	Zone int     `json:"zone,omitempty"`
}

type zoneRequest struct {
	Lon float64 `json:"lon"`
}

func forwardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}
	var req forwardRequest
	if !readJSON(w, r, &req) {
		return
	}
	result, err := utm.Forward(req.Lat, req.Lon, req.Zone)
	if err != nil {
		writeValidationOutcome(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func zoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}
	var req zoneRequest
	if !readJSON(w, r, &req) {
		return
	}
	zone, err := utm.ZoneForLongitude(req.Lon)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, zone)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, "GET")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"name": "utm-fwd", "version": "1.0.0"})
}
