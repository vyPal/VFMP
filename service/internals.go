package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

func countFiles(root string, status chan<- int) {
	var count int

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			count++
		}

		select {
		case <-ticker.C:
			status <- count
		default:
		}

		return nil
	})

	if err != nil {
		log.Printf("Failed to walk file system: %v", err)
	}

	status <- count
	close(status)
}

func walkFiles(root string, status chan<- int, trie *HybridTrie) {
	var count int

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			count++
			trie.AddPath(path)
		}

		select {
		case <-ticker.C:
			status <- count
		default:
		}

		return nil
	})

	if err != nil {
		log.Printf("Failed to walk file system: %v", err)
	}

	status <- count
	close(status)
}
