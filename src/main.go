package main

import "fmt"

import server "./atsserver"

func main() {
	fmt.Println("Initializing ATS Server...")

	s := server.CreateServer(client)
	s.Bind(":4441")
	//
	// cs := s.CandidateService
	// candidate := cs.GetTemplate()
	//
	// fmt.Println("Candidate Template:", candidate.Name)

	// 1. Listen to port 4441
	// 2. Handle telnet connection
}

func client(c chan []byte) {
	for {
		rendered := ui()
		c <- []byte(rendered)

		userChoice := <-c
		userOutput := fmt.Sprintf("You chose: %s\n", userChoice)
		c <- []byte(userOutput)
	}
}

func ui() string {
	lines :=
		`ATS Applicant Tracking System

1. Search for positions
2. Post a new position

Your choice:`

	return lines
}
