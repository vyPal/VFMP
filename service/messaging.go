package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

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

func setupIPCServer(cfg *ConfigDatabase) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		log.Fatal("Unable to listen: ", err)
	}
	defer listener.Close()

	log.Printf("IPC server listening on port %d", cfg.Server.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print("Error accepting connection: ", err)
			continue
		}

		go handleConnection(conn, cfg)
	}
}

func handleConnection(conn net.Conn, cfg *ConfigDatabase) {
	defer conn.Close()

	log.Print("New connection established")

	for {
		// Read incoming message
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Print("Error reading message: ", err)
			return
		}

		// Check message type
		var data map[string]interface{}
		err = json.Unmarshal([]byte(msg), &data)
		if err != nil {
			log.Print("Error unmarshaling message: ", err)
			return
		}

		processMessage(msg, conn, cfg)
	}
}

func processMessage(msg string, conn net.Conn, cfg *ConfigDatabase) {
	var m IPCMessage
	err := json.Unmarshal([]byte(msg), &m)
	if err != nil {
		log.Print("Error unmarshaling message: ", err)
		return
	}

	switch m.Type {
	case "count":
		var r CountRequest
		err = json.Unmarshal([]byte(m.Data), &r)
		if err != nil {
			log.Print("Error unmarshaling data: ", err)
			return
		}
		processCount(r, conn)
	case "index":
		var r IndexRequest
		err = json.Unmarshal([]byte(m.Data), &r)
		if err != nil {
			log.Print("Error unmarshaling data: ", err)
			return
		}
		processIndex(r, conn, cfg)
	case "search":
		var r SearchRequest
		err = json.Unmarshal([]byte(m.Data), &r)
		if err != nil {
			log.Print("Error unmarshaling data: ", err)
			return
		}
		processSearch(r, conn, cfg)
	case "ping":
		processPing(conn)
	case "kill":
		processKill(m, conn)
	default:
		log.Print("Unknown message type: ", m.Type)
	}
}

func processPing(conn net.Conn) {
	log.Print("Received ping message")

	// Send pong message
	pong := IPCMessage{
		Type: "pong",
	}

	pongData, err := json.Marshal(pong)
	if err != nil {
		log.Print("Error marshaling pong message: ", err)
		return
	}

	_, err = conn.Write(pongData)
	if err != nil {
		log.Print("Error writing pong message: ", err)
		return
	}
}

func processKill(req IPCMessage, conn net.Conn) {
	log.Print("Received kill message")

	if req.Data != "" {
		log.Print("Kill reason: ", req.Data)
	}

	os.Exit(0)
}

func processCount(req CountRequest, conn net.Conn) {
	log.Print("Count: ", req.Dir)

	count := make(chan int)
	go countFiles(req.Dir, count)

	// When a new value is received on the channel, send it as an json object with type "count.progress"
	for c := range count {
		msg := IPCMessage{
			Type: "count.progress",
			Data: fmt.Sprintf("%d", c),
		}
		data, err := json.Marshal(msg)
		if err != nil {
			log.Print("Error marshaling message: ", err)
			return
		}
		conn.Write(data)
		conn.Write([]byte("\n"))
	}

	// Send a message with type "count.done"
	msg := IPCMessage{
		Type: "count.done",
	}
	data, err := json.Marshal(msg)
	if err != nil {
		log.Print("Error marshaling message: ", err)
		return
	}

	conn.Write(data)
	conn.Write([]byte("\n"))
}

func processIndex(req IndexRequest, conn net.Conn, cfg *ConfigDatabase) {
	log.Print("Index: ", req.Dir)

	count := make(chan int)
	var trie HybridTrie
	log.Print("Start index")
	start := time.Now()
	go walkFiles(req.Dir, count, &trie)

	// When a new value is received on the channel, send it as an json object with type "index.progress"
	for c := range count {
		msg := IPCMessage{
			Type: "index.progress",
			Data: fmt.Sprintf("%d", c),
		}
		data, err := json.Marshal(msg)
		if err != nil {
			log.Print("Error marshaling message: ", err)
			return
		}
		conn.Write(data)
		conn.Write([]byte("\n"))
	}
	diff := time.Since(start)
	log.Print("End index. Took ", diff.Milliseconds(), "ms")

	// Send a message with type "index.done"
	msg := IPCMessage{
		Type: "index.done",
	}
	data, err := json.Marshal(msg)
	if err != nil {
		log.Print("Error marshaling message: ", err)
		return
	}

	conn.Write(data)
	conn.Write([]byte("\n"))

	trie.SaveToFile(cfg.Data.Dir + "/trie.gob")
}

func processSearch(req SearchRequest, conn net.Conn, cfg *ConfigDatabase) {
	log.Print("Search: ", req.SearchString)

	var trie HybridTrie
	start := time.Now()
	err := trie.LoadFromFile(cfg.Data.Dir + "/trie.gob")
	if err != nil {
		log.Print("Error loading trie: ", err)
		return
	}

	diff := time.Since(start)
	log.Print("Trie load took ", diff.Milliseconds(), "ms")
	start = time.Now()
	encoder := json.NewEncoder(conn)
	if req.FuzzySearch {
		res := trie.FuzzySearch(req.SearchString)
		if len(res) > req.MaxResults {
			res = res[:req.MaxResults]
		}
		diff := time.Since(start)
		log.Print("Search took ", diff.Milliseconds(), "ms")
		err := encoder.Encode(res)
		if err != nil {
			log.Printf("Error encoding fuzzy search results: %v", err)
		}
	} else {
		res := trie.Search(req.SearchString)
		if len(res) > req.MaxResults {
			res = res[:req.MaxResults]
		}
		diff := time.Since(start)
		log.Print("Search took ", diff.Milliseconds(), "ms")
		err := encoder.Encode(res)
		if err != nil {
			log.Printf("Error encoding search results: %v", err)
		}
	}
	conn.Write([]byte("\n"))
}
