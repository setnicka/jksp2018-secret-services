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
	MI5      ResultMI5
}

type Result struct {
	Completed     bool
	CompletedTime time.Time
	Tries         int
}

type ResultCIA struct {
	Result
	LastTry time.Time
}
type ResultNSA Result
type ResultKGB Result
type ResultFBI Result
type ResultPPL struct {
	Result
	PackagesTracked int
}
type ResultBIS Result
type ResultMI5 Result
