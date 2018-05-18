package main

import (
	"net/http"
	"time"

	"github.com/coreos/go-log/log"
)

const (
	MI5_LOGIN    = "EricForbes"
	MI5_PASSWORD = "Pass123!"
)

func mi5IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := getGeneralData("MI5", r)
	defer func() { executeTemplate(w, "mi5Index", data) }()

	team := server.state.GetTeam(getUser(r))
	if team != nil && team.MI5.Completed {
		http.Redirect(w, r, "/intranet", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			data.MessageType = "danger"
			data.Message = "Cannot parse login form"
			return
		}

		defer server.state.Save()
		team.MI5.Tries++

		login := r.PostFormValue("login")
		password := r.PostFormValue("password")
		log.Infof("[MI5 - %s] Trying login '%s' and password '%s'", team.Login, login, password)

		if login != MI5_LOGIN || password != MI5_PASSWORD {
			data.MessageType = "danger"
			data.Message = "Invalid credentials"
			return
		}

		log.Infof("[MI5 - %s] Completed", team.Login)
		// Everything completed
		team.MI5.Completed = true
		team.MI5.CompletedTime = time.Now()
		http.Redirect(w, r, "/intranet", http.StatusSeeOther)
	} else {

	}
}

func mi5InternalHandler(w http.ResponseWriter, r *http.Request) {
	data := getGeneralData("MI5", r)
	defer func() { executeTemplate(w, "mi5Internal", data) }()

	team := server.state.GetTeam(getUser(r))
	if team == nil || !team.MI5.Completed {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

}
