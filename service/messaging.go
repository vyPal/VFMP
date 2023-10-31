package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
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
		processSearch(r, conn)
	default:
		log.Print("Unknown message type: ", m.Type)
	}
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
}

func processIndex(req IndexRequest, conn net.Conn, cfg *ConfigDatabase) {
	log.Print("Index: ", req.Dir)

	count := make(chan int)
	var trie HybridTrie
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
	}

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

	trie.SaveToFile(cfg.Data.Dir + "/trie.gob")
}

func processSearch(req SearchRequest, conn net.Conn) {
	log.Print("Search: ", req.SearchString)
}
