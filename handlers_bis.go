package main

import (
	"net/http"
)

func bisIndexHandler(w http.ResponseWriter, r *http.Request) {
	data := getGeneralData("BIS", r)
	defer func() { executeTemplate(w, "bisIndex", data) }()

	team := server.state.GetTeam(getUser(r))
	if team != nil && team.BIS.Completed {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func bisInternalHandler(w http.ResponseWriter, r *http.Request) {
	data := getGeneralData("BIS", r)
	defer func() { executeTemplate(w, "bisInternal", data) }()

	team := server.state.GetTeam(getUser(r))
	if !team.BIS.Completed {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}
