package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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
	MaxResults   int    `json:"max" default:"10"`
}

type Match struct {
	Path    string
	Indexes []int
	Score   int
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

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter command: ")
		text, _ := reader.ReadString('\n')
		// remove the newline character
		text = strings.Replace(text, "\n", "", -1)

		// split the text into command and arguments
		parts := strings.Split(text, " ")
		command := parts[0]
		args := parts[1:]

		switch command {
		case "count":
			if len(args) != 1 {
				fmt.Println("count command requires 1 argument")
				continue
			}
			err := sendCount(args[0], conn)
			if err != nil {
				fmt.Println("Error sending count:", err)
			}
		case "index":
			if len(args) != 1 {
				fmt.Println("count command requires 1 argument")
				continue
			}
			err := sendIndex(args[0], conn)
			if err != nil {
				fmt.Println("Error sending count:", err)
			}
		case "search":
			if len(args) != 3 {
				fmt.Println("count command requires 3 arguments")
				continue
			}
			fuz, err := strconv.ParseBool(args[2])
			if err != nil {
				fmt.Println("Argument 2 is not bool")
			}
			err = sendSearch(args[0], args[1], fuz, conn)
			if err != nil {
				fmt.Println("Error sending count:", err)
			}
			// ... (your other cases)
		case "exit":
			return
		default:
			fmt.Println("Unknown command")
		}
	}
}

func sendCount(dir string, conn net.Conn) error {
	data := CountRequest{
		Dir: dir,
	}
	jsonReq, err := json.Marshal(data)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(IPCMessage{
		Type: "count",
		Data: string(jsonReq),
	})

	_, err = conn.Write(append(jsonData, []byte("\n")...))
	if err != nil {
		return err
	}
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return nil
		}

		var ipcMessage IPCMessage
		err = json.Unmarshal([]byte(message), &ipcMessage)
		if err != nil {
			fmt.Println("Error unmarshalling IPCMessage:", err)
			return nil
		}

		if ipcMessage.Type == "count.done" {
			fmt.Println("Count done")
			return nil
		} else if ipcMessage.Type == "count.progress" {
			fmt.Println("Count progress: ", ipcMessage.Data, " files")
		}
	}
}

func sendIndex(dir string, conn net.Conn) error {
	data := IndexRequest{
		Dir: dir,
	}
	jsonReq, err := json.Marshal(data)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(IPCMessage{
		Type: "index",
		Data: string(jsonReq),
	})

	_, err = conn.Write(append(jsonData, []byte("\n")...))
	if err != nil {
		return err
	}
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return nil
		}

		var ipcMessage IPCMessage
		err = json.Unmarshal([]byte(message), &ipcMessage)
		if err != nil {
			fmt.Println("Error unmarshalling IPCMessage:", err)
			return nil
		}

		if ipcMessage.Type == "index.done" {
			fmt.Println("Index done")
			return nil
		} else if ipcMessage.Type == "index.progress" {
			fmt.Println("Index progress: ", ipcMessage.Data, " files")
		}
	}
}

func sendSearch(path, search string, fuzzy bool, conn net.Conn) error {
	req := SearchRequest{
		Dir:          path,
		SearchString: search,
		FuzzySearch:  fuzzy,
		MaxResults:   10,
	}
	jsonReq, err := json.Marshal(req)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(IPCMessage{
		Type: "search",
		Data: string(jsonReq),
	})

	_, err = conn.Write(append(jsonData, '\n'))
	if err != nil {
		return err
	}

	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return nil
	}

	if req.FuzzySearch {
		var matches []Match
		err = json.NewDecoder(strings.NewReader(strings.TrimSuffix(message, "\n"))).Decode(&matches)
		if err != nil {
			fmt.Println("Error decoding matches:", err)
			return nil
		}
		fmt.Println("Matches:", matches)
	} else {
		var files []string
		err = json.NewDecoder(strings.NewReader(strings.TrimSuffix(message, "\n"))).Decode(&files)
		if err != nil {
			fmt.Println("Error decoding files:", err)
			return nil
		}
		fmt.Println("Files:", files)
	}

	return nil
}
