package main

import (
	"net/http"
	"regexp"
	"time"

	"github.com/coreos/go-log/log"
)

const (
	CIA_LOGIN = "agentfred"
)

func ciaIndexHandler(w http.ResponseWriter, r *http.Request) {
	data := getGeneralData("CIA", r)
	defer func() { executeTemplate(w, "ciaIndex", data) }()

	team := server.state.GetTeam(getUser(r))
	if team != nil && team.CIA.Completed {
		http.Redirect(w, r, "/internal", http.StatusSeeOther)
		return
	}

	//a3a.B3B3...

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			data.MessageType = "danger"
			data.Message = "Cannot parse login form"
			return
		}

		defer server.state.Save()
		team.CIA.Tries++

		login := r.PostFormValue("login")
		password := r.PostFormValue("password")
		log.Infof("[CIA - %s] Trying login '%s' and password '%s'", team.Login, login, password)

		if login != CIA_LOGIN {
			data.MessageType = "danger"
			data.Message = "Incorrect login name, only agents can log in"
			return
		} else if len(password) < 11 {
			data.MessageType = "danger"
			data.Message = "Incorrect password, password is too short"
			return
		} else if len(password) > 11 {
			data.MessageType = "danger"
			data.Message = "Incorrect password, password is too long"
			return
		}

		bpassword := []byte(password)

		reLetter := regexp.MustCompile("[a-z]")
		reBigLetter := regexp.MustCompile("[A-Z]")
		reNumbers := regexp.MustCompile("[0-9]")
		reDoubleLetters := regexp.MustCompile("[a-zA-Z][a-zA-Z]")

		smallLetters := len(reLetter.FindAll(bpassword, -1))
		bigLetters := len(reBigLetter.FindAll(bpassword, -1))
		doubleLetters := len(reDoubleLetters.FindAll(bpassword, -1))
		numbers := len(reNumbers.FindAll(bpassword, -1))
		letters := smallLetters + bigLetters
		other := len(password) - letters - numbers

		if letters < 3 {
			data.MessageType = "danger"
			data.Message = "Incorrect password, there is not enough letters"
			return
		}

		if other < 3 {
			data.MessageType = "danger"
			data.Message = "Incorrect password, there is to much letters and numbers"
			return
		}
		if doubleLetters > 0 {
			data.MessageType = "danger"
			data.Message = "Incorrect password, two letters side by side are forbidden for safety reasons"
			return
		}
		if smallLetters != bigLetters {
			data.MessageType = "danger"
			data.Message = "Incorrect password, number of lowercase letters differ from number of uppercase letters"
			return
		}
		if other < numbers {
			data.MessageType = "danger"
			data.Message = "Incorrect password, there is to much numbers"
			return
		}

		log.Infof("[CIA - %s] Completed", team.Login)
		// Everything completed
		team.CIA.Completed = true
		team.CIA.CompletedTime = time.Now()
		http.Redirect(w, r, "/internal", http.StatusSeeOther)
	} else {

	}
}

func ciaInternalHandler(w http.ResponseWriter, r *http.Request) {
	data := getGeneralData("CIA", r)
	defer func() { executeTemplate(w, "ciaInternal", data) }()

	team := server.state.GetTeam(getUser(r))
	if team == nil || !team.CIA.Completed {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}
