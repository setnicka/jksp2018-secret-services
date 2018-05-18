package main

import (
	"net/http"
	"time"

	"github.com/coreos/go-log/log"
)

const (
	BIS_LOGIN    = "reditel"
	BIS_PASSWORD = "mujmalygripen"
)

func bisIndexHandler(w http.ResponseWriter, r *http.Request) {
	data := getGeneralData("BIS", r)
	defer func() { executeTemplate(w, "bisIndex", data) }()

	team := server.state.GetTeam(getUser(r))
	if team != nil && team.BIS.Completed {
		http.Redirect(w, r, "/tajne", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			data.MessageType = "danger"
			data.Message = "Cannot parse login form"
			return
		}

		defer server.state.Save()
		team.BIS.Tries++

		login := r.PostFormValue("login")
		password := r.PostFormValue("password")
		log.Infof("[BIS - %s] Trying login '%s' and password '%s'", team.Login, login, password)

		if login != BIS_LOGIN || password != BIS_PASSWORD {
			data.MessageType = "danger"
			data.Message = "Nesprávný login, zkuste to znovu nebo kontaktujte podporu"
			return
		}

		log.Infof("[BIS - %s] Completed", team.Login)
		// Everything completed
		team.BIS.Completed = true
		team.BIS.CompletedTime = time.Now()
		http.Redirect(w, r, "/tajne", http.StatusSeeOther)
	} else {

	}
}

func bisInternalHandler(w http.ResponseWriter, r *http.Request) {
	data := getGeneralData("BIS", r)
	defer func() { executeTemplate(w, "bisInternal", data) }()

	team := server.state.GetTeam(getUser(r))
	if team == nil || !team.BIS.Completed {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}
