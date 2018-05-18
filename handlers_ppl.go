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
	PplTrackingData{"55604070459", "Horní Sytová", "512 41", "Marťanská kolonizační a.s.", "0,1 kg", "Balíček je v hajzlu. První patro, třetí dveře vpravo."},
	PplTrackingData{"16145962924", "Horní Sytová", "512 41", "Elektrárna Dukovany", "9,5 kg", "Balíček je určen na podpal."},
	PplTrackingData{"35015671123", "Horní Sytová", "512 41", "Pražský urychlovač částic", "10 tun", "TODO"},
	PplTrackingData{"44467405130", "Horní Sytová", "512 41", "Pan F, ČP", "0,19 kg", "Balíček byl úspěšně doručen. Access Granted."},
}

type pplTrackingPageData struct {
	GeneralData
	PplTrackingData
}

func pplIndexHandler(w http.ResponseWriter, r *http.Request) {
	data := getGeneralData("PPL", r)
	defer func() { executeTemplate(w, "pplIndex", data) }()

	team := server.state.GetTeam(getUser(r))
	if team != nil && team.PPL.Completed {
		http.Redirect(w, r, "/internal", http.StatusSeeOther)
		return
	}
}

func pplInternalHandler(w http.ResponseWriter, r *http.Request) {
	data := getGeneralData("PPL", r)
	defer func() { executeTemplate(w, "pplInternal", data) }()

	team := server.state.GetTeam(getUser(r))
	if team == nil || !team.PPL.Completed {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func pplTrackingHandler(w http.ResponseWriter, r *http.Request) {
	data := pplTrackingPageData{GeneralData: getGeneralData("PPL", r)}

	team := server.state.GetTeam(getUser(r))
	if team == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if team != nil && team.PPL.Completed {
		http.Redirect(w, r, "/internal", http.StatusSeeOther)
		return
	}

	defer server.state.Save()
	team.PPL.Tries++

	code := r.URL.Query().Get("zasilka_id")
	log.Infof("[PPL - %s] Trying package id: %s", team.Login, code);

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
				log.Infof("[PPL - %s] Tried to track package %d without discovering package %d", team.Login, i + 1, team.PPL.PackagesTracked + 1)
			}
		}
	}

	if team.PPL.PackagesTracked == len(packages) && !team.PPL.Completed {
		data.Message = "Access Granted"
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
