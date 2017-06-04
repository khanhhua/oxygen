package atsserver

// Candidate Candidate
type Candidate struct {
	Id          int
	Name        string
	Yob         int
	PositionIds []int
}

// Position Position
type Position struct {
	Id     int
	Name   string
	Salary int
}

// PositionQuery PositionQuery
type PositionQuery struct {
	query map[string]string
	limit int
}

// CandidateService CandidateService
type CandidateService interface {
	GetTemplate() Candidate
	Register(*Candidate) bool
}

// PositionService PositionService
type PositionService interface {
	GetTemplate() Position
	Create(*Position) bool
	Find(PositionQuery) []*Position
}
