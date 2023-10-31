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

type IPCMessage struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type CountRequest struct {
	Dir        string  `json:"dir"`
	UpdateFreq float32 `json:"ufreq" default:"10"`
}

type IndexRequest struct {
	Dir        string  `json:"dir"`
	UpdateFreq float32 `json:"ufreq" default:"10"`
}

type SearchRequest struct {
	Dir          string `json:"dir"`
	SearchString string `json:"search"`
	FuzzySearch  bool   `json:"fuzzy" default:"false"`
	MinScore     int    `json:"score" default:"0"`
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
	conn.SetDeadline(time.Now().Add(15 * time.Second))

	request := IndexRequest{
		Dir: "/home/vypal/Dokumenty",
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("Failed to marshal request: %v\n", err)
		return
	}

	message := IPCMessage{
		Type: "index",
		Data: string(jsonRequest),
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

	// Read responses from the server forever
	for {
		// Read the response from the server
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("Failed to read response from server: %v\n", err)
			return
		}

		// Print the response from the server
		fmt.Printf("Response from server: %s\n", buf[:n])
	}
}
