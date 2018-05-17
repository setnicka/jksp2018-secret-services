package state

import (
	"time"
)

type State struct {
	Teams []team
}

type team struct {
	Name     string
	Login    string
	Salt     string
	Password string
	CIA      ResultCIA
	NSA      ResultNSA
	KGB      ResultKGB
	FBI      ResultFBI
	PPL      ResultPPL
	BIS      ResultBIS
}

type Result struct {
	Completed     bool
	CompletedTime time.Time
	Tries         int
}

type ResultCIA Result
type ResultNSA Result
type ResultKGB Result
type ResultFBI Result
type ResultPPL struct {
	Result
	FirstPackage  string
	SecondPackage string
	ThirdPackage  string
}
type ResultBIS Result
