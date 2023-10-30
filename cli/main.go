package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

type ConfigDatabase struct {
	Data struct {
		Dir string `yaml:"dir" default:"/var/lib/vfmp"`
	} `yaml:"data"`
	Server struct {
		Port int `yaml:"port" default:"32768"`
	}
}

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func main() {

	// Set the path to the config file
	configFile := filepath.Join("/root", ".config", "vfmp", "config.yaml")

	// Read the config file
	var config ConfigDatabase
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	// Parse the config file
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	// Set the server address and port
	serverAddr := fmt.Sprintf("localhost:%d", config.Server.Port)

	// Connect to the server
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Printf("Failed to connect to server: %v\n", err)
		return
	}
	defer conn.Close()

	// Set a timeout for the connection
	conn.SetDeadline(time.Now().Add(10 * time.Second))

	message := Message{
		Type: "ping",
		Data: "",
	}
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Failed to marshal message: %v\n", err)
		return
	}

	// Write the json + "\n" to the server
	_, err = conn.Write(append(jsonMessage, '\n'))
	if err != nil {
		fmt.Printf("Failed to send message to server: %v\n", err)
		return
	}

	// Receive a response from the server
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Failed to receive response from server: %v\n", err)
		return
	}

	response := string(buffer[:n])
	fmt.Printf("Received response from server: %s\n", response)
}
