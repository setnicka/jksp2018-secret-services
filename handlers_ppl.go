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

var (
	PACKAGE_1 = PplTrackingData{"55604070459", "Horní Sytová", "512 41", "TODO Odesílatel", "0,1 kg", "Balíček je v hajzlu. První patro, třetí dveře vpravo."} // TODO: Pozice
	PACKAGE_2 = PplTrackingData{"16145962924", "Horní Sytová", "512 41", "TODO Odesílatel", "0,15 kg", "Balíček je určen na podpal."} // TODO: Pozice
	PACKAGE_3 = PplTrackingData{"44467405130", "Horní Sytová", "512 41", "TODO Odesílatel", "0,19 kg", "Balíček je neumístěn."} // TODO: Pozice
)

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

	if code == PACKAGE_1.Code {
		if (team.PPL.PackagesTracked == 0) {
			team.PPL.PackagesTracked = 1
			log.Infof("[PPL - %s] Tracked package 1", team.Login)
		}
		if (team.PPL.PackagesTracked >= 1) {
			data.PplTrackingData = PACKAGE_1
		}
	} else if code == PACKAGE_2.Code {
		if (team.PPL.PackagesTracked == 1) {
			team.PPL.PackagesTracked = 2
			log.Infof("[PPL - %s] Tracked package 2", team.Login)
		}
		if (team.PPL.PackagesTracked >= 2) {
			data.PplTrackingData = PACKAGE_2
		} else {
			log.Infof("[PPL - %s] Tried to track package 2 without discovering package 1", team.Login)
		}
	} else if code == PACKAGE_3.Code {
		if (team.PPL.PackagesTracked == 2) {
			team.PPL.PackagesTracked = 3
			log.Infof("[PPL - %s] Tracked package 3", team.Login)

			team.PPL.Completed = true
			team.PPL.CompletedTime = time.Now()
		}
		if (team.PPL.PackagesTracked >= 3) {
			data.PplTrackingData = PACKAGE_3
		} else {
			log.Infof("[PPL - %s] Tried to track package 3 without discovering package 2", team.Login)
		}
	}
	if (data.PplTrackingData.Code != "") {
		defer func() { executeTemplate(w, "pplTracking", data) }()
	} else {
	defer func() { executeTemplate(w, "pplPackageNotFound", data) }()
	}
}
