package main

import (
	"net/http"
)

func nsaIndexHandler(w http.ResponseWriter, r *http.Request) {
	data := getGeneralData("NSA", r)
	defer func() { executeTemplate(w, "nsaIndex", data) }()

	team := server.state.GetTeam(getUser(r))
	if team != nil && team.NSA.Completed {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func nsaInternalHandler(w http.ResponseWriter, r *http.Request) {
	data := getGeneralData("NSA", r)
	defer func() { executeTemplate(w, "nsaInternal", data) }()

	team := server.state.GetTeam(getUser(r))
	if team == nil || !team.NSA.Completed {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}
