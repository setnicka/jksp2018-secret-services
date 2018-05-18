package main

import (
	"net/http"
	"regexp"
	"time"
	"unicode"

	"github.com/coreos/go-log/log"
)

const (
	CIA_LOGIN  = "agentfred"
	TIME_LIMIT = time.Second * 60
)

func ciaIndexHandler(w http.ResponseWriter, r *http.Request) {
	data := getGeneralData("CIA", r)
	defer func() { executeTemplate(w, "ciaIndex", data) }()

	team := server.state.GetTeam(getUser(r))
	if team != nil && team.CIA.Completed {
		http.Redirect(w, r, "/internal", http.StatusSeeOther)
		return
	}

	// e.F.g.H.i.J

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			data.MessageType = "danger"
			data.Message = "Cannot parse login form"
			return
		}

		if time.Since(team.CIA.LastTry) < TIME_LIMIT {
			data.MessageType = "danger"
			data.Message = "Minimal time between two login attempts is 60 seconds"
			return
		}
		defer server.state.Save()
		team.CIA.Tries++
		team.CIA.LastTry = time.Now()

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
			data.Message = "Incorrect password, there isn't enough letters"
			return
		}

		if other < 3 {
			data.MessageType = "danger"
			data.Message = "Incorrect password, there are too many letters and numbers"
			return
		}
		if doubleLetters > 0 {
			data.MessageType = "danger"
			data.Message = "Incorrect password, two letters side by side are forbidden for safety reasons"
			return
		}
		if smallLetters != bigLetters {
			data.MessageType = "danger"
			data.Message = "Incorrect password, the number of lowercase letters differs from the number of uppercase letters"
			return
		}
		if other < numbers {
			data.MessageType = "danger"
			data.Message = "Incorrect password, there is too many numbers"
			return
		}

		// pismena nesmi byt stejna
		wasLetter := map[rune]bool{}
		for _, r := range password {
			if unicode.IsLetter(r) {
				if _, found := wasLetter[r]; found {
					data.MessageType = "danger"
					data.Message = "Incorrect password, there cannot be two same letters"
					return
				}
				wasLetter[r] = true
			}
		}

		// pismena nejsou v abecednim poradi
		lastLetter := 'a'
		for _, r := range password {
			if unicode.IsLetter(r) {
				r = unicode.ToLower(r)
				if r < lastLetter {
					data.MessageType = "danger"
					data.Message = "Incorrect password, letters must be in alphabetical order"
					return
				}
				lastLetter = r
			}
		}

		// na tretim miste neni f
		if password[2] != 'f' {
			data.MessageType = "danger"
			data.Message = "Incorrect password, there isn't the first letter of your first name on the third position in the password"
			return
		}

		reAroundNumber := regexp.MustCompile("[a-zA-Z][0-9][a-zA-Z]")
		reAroundNumberSpecial := regexp.MustCompile("[^0-9a-zA-Z][0-9][^0-9a-zA-Z]")

		// cislo nesmi mit z obou stran specialni znak
		if reAroundNumberSpecial.Match(bpassword) {
			data.MessageType = "danger"
			data.Message = "Incorrect password, number cannot be surrounded by special characters from both sides"
			return
		}

		// cislo nesmi mit z obou stran pismeno
		if reAroundNumber.Match(bpassword) {
			data.MessageType = "danger"
			data.Message = "Incorrect password, number cannot be surrounded by letters from both sides"
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
