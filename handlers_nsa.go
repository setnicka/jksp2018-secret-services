package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-log/log"
)

const (
	NSA_LOGIN    = "patterson"
	NSA_PASSWORD = "laoo,rpe" // "password" on Dvorak keyboard
)

func nsaIndexHandler(w http.ResponseWriter, r *http.Request) {
	data := getGeneralData("NSA", r)
	defer func() { executeTemplate(w, "nsaIndex", data) }()

	team := server.state.GetTeam(getUser(r))
	if team != nil && team.NSA.Completed {
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
		team.NSA.Tries++

		login := r.PostFormValue("login")
		password := r.PostFormValue("password")
		log.Infof("[NSA - %s] Trying login '%s' and password '%s'", team.Login, login, password)

		if strings.ToLower(login) == NSA_LOGIN && password == NSA_PASSWORD {
			log.Infof("[NSA - %s] Completed", team.Login)
			// Everything completed
			team.NSA.Completed = true
			team.NSA.CompletedTime = time.Now()
			http.Redirect(w, r, "/intranet", http.StatusSeeOther)
		} else {
			data.MessageType = "danger"
			data.Message = "Incorrect login or password, check your spelling"
		}
	}
}

func nsaIntranetHandler(w http.ResponseWriter, r *http.Request) {
	data := getGeneralData("NSA", r)
	defer func() { executeTemplate(w, "nsaIntranet", data) }()

	team := server.state.GetTeam(getUser(r))
	if team == nil || !team.NSA.Completed {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}
