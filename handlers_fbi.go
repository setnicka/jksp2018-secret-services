package main

import (
	"net/http"
	"time"

	"github.com/coreos/go-log/log"
)

const (
	FBI_LOGIN    = "agentjones"
	FBI_PASSWORD = "restart"
)

func fbiIndexHandler(w http.ResponseWriter, r *http.Request) {
	data := getGeneralData("FBI", r)
	defer func() { executeTemplate(w, "fbiIndex", data) }()

	team := server.state.GetTeam(getUser(r))
	if team != nil && team.FBI.Completed {
		http.Redirect(w, r, "/private", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			data.MessageType = "danger"
			data.Message = "Cannot parse login form"
			return
		}

		defer server.state.Save()
		team.FBI.Tries++

		login := r.PostFormValue("login")
		password := r.PostFormValue("pass")
		log.Infof("[FBI - %s] Trying login '%s' and password '%s'", team.Login, login, password)

		if login != FBI_LOGIN {
			data.MessageType = "danger"
			data.Message = "Incorrect login name, only agents can log in"
			return
		}

		if password != FBI_PASSWORD {
			data.MessageType = "danger"
			data.Message = "Incorrect password"
			return
		}

		log.Infof("[FBI - %s] Completed", team.Login)
		// Everything completed
		team.FBI.Completed = true
		team.FBI.CompletedTime = time.Now()
		http.Redirect(w, r, "/private", http.StatusSeeOther)
	} else {

	}
}

func fbiInternalHandler(w http.ResponseWriter, r *http.Request) {
	data := getGeneralData("FBI", r)
	defer func() { executeTemplate(w, "fbiIndex", data) }()

	team := server.state.GetTeam(getUser(r))
	if team == nil || !team.FBI.Completed {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data.MessageType = "success"
	data.Message = "Access Granted"
}
