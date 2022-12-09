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


func broadcastMessage(self string, connectedClients *map[string]*Client, message string) {
	for _, client := range *connectedClients {
		if client.name == self {
			continue // Do not send the message back to the client who sent it
		}
		client.connection.Write([]byte(message))
	}
}

func clientHandler(client *Client, connectedClients *map[string]*Client) {
	// Send a welcome message to the client
	client.connection.Write([]byte("Welcome to the chat server, " + client.name + "!\n"))
	scanner := bufio.NewScanner(client.connection) // Create a scanner to read the client's input
	for scanner.Scan() {
		msgSent := scanner.Text() // Get the message sent by the client
		// Broadcast the message to all other clients
		formattedMessage := "> " + client.name + ": " + msgSent + "\n"
		fmt.Print(formattedMessage)
		broadcastMessage(client.name, connectedClients, formattedMessage)
	}

	// If the client has disconnected, remove them from the map
	disconnectedName := client.name
	delete(*connectedClients, client.name)
	fmt.Println("Client disconnected")
	broadcastMessage(disconnectedName, connectedClients, "\n"+disconnectedName+" has disconnected.\n")
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
		connection.Write([]byte("Please enter your name: "))
		reader := bufio.NewReader(connection)
		name, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected before entering name")
			continue
		}
		name = name[:len(name)-1] // Removes the newline character in name
		// Check if the name is taken by doing a map lookup
		existingClient := conenctedClients[name]
		if existingClient != nil {
			// The name is taken, ask for a new name
			connection.Write([]byte("Name is taken. Please choose another name."))
			fmt.Println("Closing connection")
			connection.Close()
			continue
		}
		newClient := Client{connection, name} // Create a new client
		conenctedClients[name] = &newClient // Add to the map
		broadcastMessage(name, &conenctedClients, name+" has connected.\n")
		go clientHandler(&newClient, &conenctedClients)
	}
}
