package main

import (
	"net/http"
	"time"

	"github.com/coreos/go-log/log"
)

const (
	KGB_LOGIN = "igorkarkarov"
	KGB_PASSWORD = "MarxRulles"
)

func kgbIndexHandler(w http.ResponseWriter, r *http.Request) {
	data := getGeneralData("KGB", r)
	defer func() { executeTemplate(w, "kgbIndex", data) }()

	team := server.state.GetTeam(getUser(r))
	if team != nil && team.KGB.Completed {
		http.Redirect(w, r, "/интранет", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			data.MessageType = "danger"
			data.Message = "Cannot parse login form"
			return
		}

		defer server.state.Save()
		team.KGB.Tries++

		a1 := r.PostFormValue("a1") // první půlka hesla
		a2 := r.PostFormValue("a2") // login bez prvního písmenka
		a3 := r.PostFormValue("a3") // login
		a4 := r.PostFormValue("a4") // prázdné
		a5 := r.PostFormValue("a5") // 2017
		a6 := r.PostFormValue("a6") // druhý znak hesla
		a7 := r.PostFormValue("a7") // password
		a8 := r.PostFormValue("a8") // cucoriedka
		a9 := r.PostFormValue("a9") // login pozpátku
		a10 := r.PostFormValue("a10") // každý druhý znak hesla
		log.Infof("[KGB - %s] Trying login '%s' and password '%s'", team.Login, a3, a7)
		pw_half := KGB_PASSWORD[:len(KGB_PASSWORD) / 2]
		login_wo_first := KGB_LOGIN[1:]
		login_reversed := ""
		for i := len(KGB_LOGIN) - 1; i >= 0; i -= 1 {
			login_reversed += string(KGB_LOGIN[i])
		}

		pw_every_second := ""
		for i := 0; i < len(KGB_PASSWORD); i += 2 {
			pw_every_second += string(KGB_PASSWORD[i])
		}

		fail := a1 != pw_half ||
		     a2 != login_wo_first ||
		     a3 != KGB_LOGIN ||
		     len(a4) != 0 ||
		     a5 != "2017" ||
		     len(a6) < 1 || a6[0] != KGB_PASSWORD[1] ||
		     a7 != KGB_PASSWORD ||
		     a8 != "cucoriedka" ||
		     a9 != login_reversed ||
		     a10 != pw_every_second

		if fail {
			data.MessageType = "danger"
			data.Message = "Пассворда нет валидна"
			return
		}

		log.Infof("[KGB - %s] Completed", team.Login)
		// Everything completed
		team.KGB.Completed = true
		team.KGB.CompletedTime = time.Now()
		http.Redirect(w, r, "/internal", http.StatusSeeOther)
	} else {

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
