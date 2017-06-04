package main

import "fmt"
import "strconv"

import server "./atsserver"

func main() {
	fmt.Println("Initializing ATS Server...")

	s := server.CreateServer(clientDelegate)
	s.Bind(":4441")
	//
	// cs := s.CandidateService
	// candidate := cs.GetTemplate()
	//
	// fmt.Println("Candidate Template:", candidate.Name)

	// 1. Listen to port 4441
	// 2. Handle telnet connection
}

func clientDelegate(client server.Client) {
	var userChoice string

	for {
		// Show UI to client
		fmt.Println("clientDelegate: client.Screen:", client.Screen)
		rendered := ui(client)
		client.Send(rendered)

		switch client.Screen {
		case server.ScreenMainMenu:
			// Get user input
			userChoice = client.Receive()
			fmt.Println("clientDelegate: You chose:", userChoice)

			switch userChoice {
			case "1":
				client.Screen = server.ScreenSearchPositions
			case "2":
				client.Screen = server.ScreenPostPosition
			}

		case server.ScreenPostPosition:
			client.Send("Position name: ")
			positionName := client.Receive()
			client.Send("Position salary: ")
			positionSalary := client.Receive()

			fmt.Println("Preparing to post a new position", positionName, positionSalary)
			iPositionSalary, _ := strconv.Atoi(positionSalary)
			result := client.CreatePosition(positionName, iPositionSalary)

			if result == true {
				client.Send("A new position has been created\n")
			}
			client.Send("\n\n")

			client.Screen = server.ScreenMainMenu

		case server.ScreenSearchPositions:
			client.Send("Min salary: ")
			minSalary := client.Receive()
			client.Send("Position name: ")
			positionName := client.Receive()

			fmt.Println("Preparing to search for positions")
			query := make(map[string]string)
			query["minSalary"] = minSalary
			query["positionName"] = positionName

			positions := client.SearchPositions(query)

			client.Send("==================================\n")
			client.Send("Positions found\n")
			client.Send("----------------------------------\n")

			for _, position := range positions {
				client.Send(fmt.Sprintf("%d. %s\t\t%d\n", position.Id, position.Name, position.Salary))
			}
			client.Send("----------------------------------\n\n")

			client.Send("0. Back to main menu\n")
			client.Send("1. Apply for a position\n\n")
			client.Send("X. Quit\n")
			client.Send("Your choice: ")

			userChoice = client.Receive()
			switch userChoice {
			case "0":
				client.Screen = server.ScreenMainMenu
			case "1":
				client.Screen = server.ScreenApplyForPosition
			}

		case server.ScreenApplyForPosition:
			client.Send("Position ID: ")
			positionId := client.Receive()
			client.Send("Your name: ")
			applicantName := client.Receive()
			client.Send("Your year of birth: ")
			applicantYob := client.Receive()

			iPositionId, _ := strconv.Atoi(positionId)
			iApplicantYob, _ := strconv.Atoi(applicantYob)

			result := client.ApplyPosition(iPositionId, applicantName, iApplicantYob)
			if result == true {
				client.Send(fmt.Sprintf("Candidate %s has applied for a job\n", applicantName))
				client.Screen = server.ScreenMainMenu
			}
		}

		// For logging purpose only
		// userOutput := fmt.Sprintf("clientDelegate: You chose: %s\n", userChoice)
		// client.Send(userOutput)
	}

	client.Quit()
}

func ui(client server.Client) string {
	screen := client.Screen

	fmt.Println("Rendering UI for screen#", screen)

	switch screen {
	case server.ScreenMainMenu:
		return `ATS Applicant Tracking System
-----------------------------
1. Search for positions
2. Post a new position

X. Quit

Your choice:`
	case server.ScreenSearchPositions:
		return `ATS Applicant Tracking System
-----------------------------
1. Search for positions
=============================
Specify your criteria...
`
	case server.ScreenPostPosition:
		return `ATS Applicant Tracking System
-----------------------------
2. Post a new position
=============================
`
	case server.ScreenApplyForPosition:
		return `ATS Applicant Tracking System
-----------------------------
1. Post a new position
  1. Apply for a position
=============================
`
	default:
		return "Unknown"
	}
}
