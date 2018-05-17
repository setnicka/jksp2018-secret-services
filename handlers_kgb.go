package main

import (
	"net/http"
)

func kgbIndexHandler(w http.ResponseWriter, r *http.Request) {
	data := getGeneralData("KGB", r)
	defer func() { executeTemplate(w, "kgbIndex", data) }()

	team := server.state.GetTeam(getUser(r))
	if team != nil && team.KGB.Completed {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func kgbInternalHandler(w http.ResponseWriter, r *http.Request) {
	data := getGeneralData("KGB", r)
	defer func() { executeTemplate(w, "kgbInternal", data) }()

	team := server.state.GetTeam(getUser(r))
	if team == nil || !team.KGB.Completed {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}
