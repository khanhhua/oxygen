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
	for {
		// Show UI to client
		fmt.Println("clientDelegate: client.Screen:", client.Screen)
		rendered := ui(client.Screen)
		client.Send(rendered)

		if client.Screen == server.ScreenMainMenu {
			// Get user input
			userChoice := client.Receive()
			fmt.Println("clientDelegate: You chose:", userChoice)

			switch userChoice {
			case "1":
				client.Screen = server.ScreenSearchPositions
			case "2":
				client.Screen = server.ScreenPostPosition
			}
		} else if client.Screen == server.ScreenSearchPositions {
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

			for i := 0; i < len(positions); i++ {
				client.Send(fmt.Sprintf("%d. %s\t\t%d\n", i+1, positions[i].Name, positions[i].Salary))
			}
			client.Send("----------------------------------\n\n")

			client.Screen = server.ScreenMainMenu
		} else if client.Screen == server.ScreenPostPosition {
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
		}

		// For logging purpose only
		// userOutput := fmt.Sprintf("clientDelegate: You chose: %s\n", userChoice)
		// client.Send(userOutput)
	}
}

func ui(screen int) string {
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
	default:
		return "Unknown"
	}
}
