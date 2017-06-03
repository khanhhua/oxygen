package atsserver

type Candidate struct {
	Id   int
	Name string
	Yob  int
}

type Position struct {
	Id     int
	Name   string
	Salary int
}

type CandidateService interface {
	GetTemplate() *Candidate
	Register(Candidate) bool
}

type PositionService interface {
	GetTemplate() *Position
	Create(Position) bool
}
