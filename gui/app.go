package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
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

// App struct
type App struct {
	ctx  context.Context
	conn net.Conn
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
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

	a.conn = conn
}

func (a *App) shutdown(ctx context.Context) {
	a.conn.Close()
}

func (a *App) SendCount(dir string) error {
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

	_, err = a.conn.Write(append(jsonData, []byte("\n")...))
	if err != nil {
		return err
	}
	reader := bufio.NewReader(a.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			runtime.EventsEmit(a.ctx, "console-output", "Error reading from connection:", err)
			return nil
		}

		var ipcMessage IPCMessage
		err = json.Unmarshal([]byte(message), &ipcMessage)
		if err != nil {
			runtime.EventsEmit(a.ctx, "console-output", "Error unmarshalling IPCMessage:", err)
			return nil
		}

		if ipcMessage.Type == "count.done" {
			runtime.EventsEmit(a.ctx, "count.done")
			return nil
		} else if ipcMessage.Type == "count.progress" {
			runtime.EventsEmit(a.ctx, "count.progress", ipcMessage.Data)
		}
	}
}

func (a *App) SendIndex(dir string) error {
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

	_, err = a.conn.Write(append(jsonData, []byte("\n")...))
	if err != nil {
		return err
	}
	reader := bufio.NewReader(a.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			runtime.EventsEmit(a.ctx, "console-output", "Error reading from connection: "+err.Error())
			return nil
		}

		var ipcMessage IPCMessage
		err = json.Unmarshal([]byte(message), &ipcMessage)
		if err != nil {
			runtime.EventsEmit(a.ctx, "console-output", "Error unmarshalling IPCMessage: "+err.Error())
			return nil
		}

		if ipcMessage.Type == "index.done" {
			runtime.EventsEmit(a.ctx, "index.done")
			return nil
		} else if ipcMessage.Type == "index.progress" {
			runtime.EventsEmit(a.ctx, "index.progress", ipcMessage.Data)
		}
	}
}

func (a *App) SendSearch(path, search string, fuzzy bool) error {
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

	_, err = a.conn.Write(append(jsonData, '\n'))
	if err != nil {
		return err
	}

	reader := bufio.NewReader(a.conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		runtime.EventsEmit(a.ctx, "console-output", "Error reading from connection: "+err.Error())
		return nil
	}

	if req.FuzzySearch {
		var matches []Match
		err = json.NewDecoder(strings.NewReader(strings.TrimSuffix(message, "\n"))).Decode(&matches)
		if err != nil {
			runtime.EventsEmit(a.ctx, "console-output", "Error decoding matches: "+err.Error())
			return nil
		}
		runtime.EventsEmit(a.ctx, "seatch.results.fuzzy", matches)
	} else {
		var files []string
		err = json.NewDecoder(strings.NewReader(strings.TrimSuffix(message, "\n"))).Decode(&files)
		if err != nil {
			runtime.EventsEmit(a.ctx, "console-output", "Error decoding files: "+err.Error())
			return nil
		}
		runtime.EventsEmit(a.ctx, "search.results", files)
	}

	return nil
}
