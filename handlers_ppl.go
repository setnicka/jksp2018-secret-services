package main

import (
	"net/http"
	"time"

	"github.com/coreos/go-log/log"
)

type PplTrackingData struct {
	Code string
	City string
	ZipCode string
	Customer string
	Weight string
	Location string
}

var packages = []PplTrackingData{
	PplTrackingData{"55604070459", "Horní Sytová", "512 41", "TODO Odesílatel", "0,1 kg", "Balíček je v hajzlu. První patro, třetí dveře vpravo."},
	PplTrackingData{"16145962924", "Horní Sytová", "512 41", "TODO Odesílatel", "0,15 kg", "Balíček je určen na podpal."},
	PplTrackingData{"44467405130", "Horní Sytová", "512 41", "TODO Odesílatel", "0,19 kg", "Balíček je neumístěn."},
}

type pplTrackingPageData struct {
	GeneralData
	PplTrackingData
}

func pplIndexHandler(w http.ResponseWriter, r *http.Request) {
	data := getGeneralData("PPL", r)
	defer func() { executeTemplate(w, "pplIndex", data) }()
}

func pplTrackingHandler(w http.ResponseWriter, r *http.Request) {
	data := pplTrackingPageData{GeneralData: getGeneralData("PPL", r)}

	team := server.state.GetTeam(getUser(r))
	if team == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	defer server.state.Save()
	team.PPL.Tries++

	code := r.URL.Query().Get("zasilka_id")
	log.Debugf("Zasilka id: %s", code);

	for i, pack := range packages {
		if code == pack.Code {
			if team.PPL.PackagesTracked == i {
				team.PPL.PackagesTracked = i + 1
				log.Infof("[PPL - %s] Tracked package %i", team.Login, i + 1)
			}
			if (team.PPL.PackagesTracked >= i + 1) {
				data.PplTrackingData = pack
			} else {
				data.Message = "Zásilka nebyla nalezena kvůli chybě v časoprostorovém kontinuu."
				log.Infof("[PPL - %s] Tried to track package %i without discovering package %i", team.Login, i + 1, team.PPL.PackagesTracked + 1)
			}
		}
	}
	if team.PPL.PackagesTracked == 3 && !team.PPL.Completed {
			team.PPL.Completed = true
			team.PPL.CompletedTime = time.Now()
	}

	if data.PplTrackingData.Code != "" {
		defer func() { executeTemplate(w, "pplTracking", data) }()
	} else {
		if data.Message == "" {
			data.Message = "Zásilka nebyla nalezena"
		}
		defer func() { executeTemplate(w, "pplPackageNotFound", data) }()
	}
}
