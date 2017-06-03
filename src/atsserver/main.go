package atsserver

import (
	"bufio"
	"fmt"
	"net"
)

// CreateServer Create an instance of ATS server which can bind to a tcp/ip ports
// and handle telnet connection
func CreateServer(client ClientDelegate) Server {
	fmt.Println("Creating ATS Server...")

	cs := CandidateServiceImpl{}
	ps := PositionServiceImpl{}

	s := Server{CandidateService: cs, PositionService: ps, client: client}

	return s
}

// Server ATS server
type Server struct {
	CandidateService CandidateService
	PositionService  PositionService
	client           ClientDelegate
}

// Bind Allow ATS to bind to a tcp/ip port
func (s Server) Bind(port string) bool {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println("ATS Server is listening to port", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		// Let's go with single client model
		handleConnection(conn, s.client)
	}
}

func handleConnection(conn net.Conn, client ClientDelegate) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	channel := make(chan []byte)
	go client(channel)

	for {
		// Maintain a conversation with client in turns
		// First speaker is server, saying hello
		userOutput := <-channel
		fmt.Println("Replying to client...")
		conn.Write(userOutput)

		userInput, _, err := reader.ReadLine()
		fmt.Println("Received from client:", userInput)

		if err != nil {
			fmt.Println("ERROR: An error has occurred while handling client connection!")
			conn.Close()
			break
		}
		channel <- userInput
	}

	conn.Close()
	return
}

// ClientDelegate Client handler for telnet connection
type ClientDelegate func(chan []byte)

// CandidateServiceImpl Implementation of CandidateService interface
type CandidateServiceImpl struct{}

// GetTemplate Implement CandidateService:GetTemplate
func (cs CandidateServiceImpl) GetTemplate() *Candidate {
	template := Candidate{Name: "Default", Yob: 2000, Id: 0}

	return &template
}

// Register Implement CandidateService:Register
func (cs CandidateServiceImpl) Register(Candidate) bool {
	fmt.Println("CandidateServiceImpl.register(Candidate)")

	return true
}

// PositionServiceImpl Implementation of PositionService
type PositionServiceImpl struct{}

// GetTemplate Implement PositionService:PositionService
func (ps PositionServiceImpl) GetTemplate() *Position {
	template := Position{Name: "Default", Salary: 20000, Id: 0}

	return &template
}

// Create Implement PositionService:Create
func (ps PositionServiceImpl) Create(Position) bool {
	fmt.Println("PositionServiceImpl.create(Position)")

	return true
}
