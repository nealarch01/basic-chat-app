package main

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	connection net.Conn
	name       string
}

func clientHandler(client *Client, connectedClients map[string]*Client) {
	// Send a welcome message to the client
	client.connection.Write([]byte("Welcome to the chat server, " + client.name + "!\n"))
	scanner := bufio.NewScanner(client.connection) // Create a scanner to read the client's input
	for scanner.Scan() {
		msgSent := scanner.Text() // Get the message sent by the client
		fmt.Println(msgSent)
		// Broadcast the message to all other clients
		for _, otherClient := range connectedClients {
			if otherClient.name == client.name {
				continue
			}
			chatMessage := "> " + client.name + ": " + msgSent + "\n"
			otherClient.connection.Write([]byte(chatMessage))
		}
	}

	// If the client has disconnected, remove them from the map
	delete(connectedClients, client.name)
	fmt.Println("Client disconnected")
	// fmt.Println(connectedClients[client.name]) // Should print <nil> if the client has been removed
	client.connection.Close()
}

func main() {
	conenctedClients := make(map[string]*Client)
	listener, err := net.Listen("tcp", ":5440")
	if err != nil { // An error may occur if the port is already in use
		fmt.Println("Error listening:", err.Error())
	}
	fmt.Println("Listening on port 5440")

	// Accept incoming requests
	for {
		connection, err := listener.Accept() // Accept the new connection; is a blocking call
		fmt.Println("New client connected")
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		// Print connection details
		connection.Write([]byte("Please enter your name: "))
		reader := bufio.NewReader(connection)
		name, _ := reader.ReadString('\n')
		// Remove the new line in name
		name = name[:len(name)-1]
		// Check if the name is taken by doing a map lookup
		existingClient := conenctedClients[name]
		if existingClient != nil {
			// The name is taken, ask for a new name
			connection.Write([]byte("Name is taken. Please choose another name."))
			fmt.Println("Closing connection")
			connection.Close()
			continue
		}
		// Add the client to the map
		newClient := Client{connection, name}
		conenctedClients[name] = &newClient
		go clientHandler(&newClient, conenctedClients)
	}
}
