package atsserver

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// CreateServer Create an instance of ATS server which can bind to a tcp/ip ports
// and handle telnet connection
func CreateServer(clientDelegate ClientDelegate) Server {
	fmt.Println("Creating ATS Server...")

	cs := CandidateServiceImpl{}
	ps := PositionServiceImpl{}

	s := Server{CandidateService: cs, PositionService: &ps, clientDelegate: clientDelegate}

	return s
}

// ScreenMainMenu 100: Main menu
const (
	ScreenMainMenu int = iota
	ScreenSearchPositions
	ScreenPostPosition
)

// Server ATS server
type Server struct {
	CandidateService CandidateService
	PositionService  PositionService
	clientDelegate   ClientDelegate
}

// Client Client
type Client struct {
	In     chan string
	Out    chan string
	Screen int
	server Server
}

// Send Send
func (c Client) Send(s string) {
	c.Out <- s
}

// Receive Receive
func (c Client) Receive() string {
	s := <-c.In

	return s
}

// CreatePosition CreatePosition
func (c Client) CreatePosition(name string, salary int) bool {
	server := c.server

	position := server.PositionService.GetTemplate()
	position.Name = name
	position.Salary = salary

	return server.PositionService.Create(&position)
}

// SearchPositions SearchPositions
func (c Client) SearchPositions(query map[string]string) []*Position {
	server := c.server

	positionQuery := PositionQuery{query, 10}
	result := server.PositionService.Find(positionQuery)

	return result
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
		s.handleConnection(conn)
	}
}

func (s Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	in := make(chan string)
	out := make(chan string)
	client := Client{In: in, Out: out, Screen: ScreenMainMenu, server: s}

	go s.clientDelegate(client)

	go func(out <-chan string) {
		for {
			userOutput := <-out
			fmt.Println("Replying to client...")
			conn.Write([]byte(userOutput))
		}
	}(out)

	for {
		// Maintain a conversation with client in turns
		// First speaker is server, saying hello
		userInputB, _, err := reader.ReadLine()

		if err != nil {
			fmt.Println("ERROR: An error has occurred while handling client connection!")
			conn.Close()
			break
		}

		userInput := string(userInputB)
		fmt.Println("Received from client:", userInput)

		if userInput == "X" {
			fmt.Println("WARNING: Client quit")
			conn.Close()
		}

		in <- userInput
	}

	close(in)
	close(out)
	conn.Close()
	return
}

// ClientDelegate Client handler for telnet connection
type ClientDelegate func(Client)

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
type PositionServiceImpl struct {
	positions []*Position
}

// GetTemplate Implement PositionService:PositionService
func (ps *PositionServiceImpl) GetTemplate() Position {
	template := Position{Name: "Default", Salary: 20000, Id: 0}

	return template
}

// Create Implement PositionService:Create
func (ps *PositionServiceImpl) Create(p *Position) bool {
	fmt.Println("PositionServiceImpl.create(Position)")

	fmt.Println("==================================")
	fmt.Println("Current positions")
	fmt.Println("----------------------------------")
	for i := 0; i < len(ps.positions); i++ {
		fmt.Printf("%d. %s\t\t%d\n", i+1, ps.positions[i].Name, ps.positions[i].Salary)
	}
	fmt.Println("==================================")

	ps.positions = append(ps.positions, p)

	fmt.Println("==================================")
	fmt.Println("New positions")
	fmt.Println("----------------------------------")
	for i := 0; i < len(ps.positions); i++ {
		fmt.Printf("%d. %s\t\t%d\n", i+1, ps.positions[i].Name, ps.positions[i].Salary)
	}
	fmt.Println("==================================")

	return true
}

// Find Find
func (ps *PositionServiceImpl) Find(query PositionQuery) []*Position {
	var result = make([]*Position, 0)

	iMinSalary, error := strconv.Atoi(query.query["minSalary"])
	if error != nil {
		return result
	}

	positionName := query.query["positionName"]

	for _, item := range ps.positions {
		if (item.Salary >= iMinSalary) && (strings.Index(item.Name, positionName) == 0) {
			result = append(result, item)
		}
	}

	return result
}
