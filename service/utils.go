package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"syscall"
	"time"
)

func YesNoPrompt(label string, def bool) bool {
	choices := "Y/n"
	if !def {
		choices = "y/N"
	}

	r := bufio.NewReader(os.Stdin)
	var s string

	for {
		fmt.Fprintf(os.Stderr, "%s (%s) ", label, choices)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "" {
			return def
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}

func tryKillDaemon(port int) {
	serverAddr := fmt.Sprintf("localhost:%d", port)

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Printf("Failed to connect to server: %v\n", err)
		return
	}
	defer conn.Close()

	// Set a timeout for the connection
	conn.SetDeadline(time.Now().Add(10 * time.Second))

	// Send kill message
	message := IPCMessage{
		Type: "kill",
		Data: "Daemon restart",
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

	// Relaunch the current process
	args := os.Args
	executable, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}

	err = syscall.Exec(executable, args, os.Environ())
	if err != nil {
		log.Fatalf("Failed to relaunch process: %v", err)
	}
}
