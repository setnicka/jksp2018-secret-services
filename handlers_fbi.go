package main

import (
	"net/http"
)

func fbiIndexHandler(w http.ResponseWriter, r *http.Request) {
	data := getGeneralData("FBI", r)
	defer func() { executeTemplate(w, "fbiIndex", data) }()

	team := server.state.GetTeam(getUser(r))
	if team != nil && team.FBI.Completed {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func fbiInternalHandler(w http.ResponseWriter, r *http.Request) {
	data := getGeneralData("FBI", r)
	defer func() { executeTemplate(w, "fbiInternal", data) }()

	team := server.state.GetTeam(getUser(r))
	if team == nil || !team.FBI.Completed {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}
