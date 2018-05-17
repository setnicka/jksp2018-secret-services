package main

import (
	"net/http"
)

func pplIndexHandler(w http.ResponseWriter, r *http.Request) {
	data := getGeneralData("PPL", r)
	defer func() { executeTemplate(w, "pplIndex", data) }()

	team := server.state.GetTeam(getUser(r))
	if team != nil && team.PPL.Completed {
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
