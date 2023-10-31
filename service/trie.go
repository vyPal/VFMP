package main

import (
	"compress/gzip"
	"encoding/gob"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/sahilm/fuzzy"
)

type TrieNode struct {
	IsEndOfWord bool
	Children    map[string]*TrieNode
}

type HybridTrie struct {
	Root TrieNode
}

func (t *HybridTrie) AddPath(path string) {
	path = strings.ReplaceAll(path, "\\", "/") // Normalize path
	parts := strings.Split(path, "/")          // Split path into parts
	node := &t.Root
	for _, part := range parts {
		if _, ok := node.Children[part]; !ok {
			if node.Children == nil {
				node.Children = make(map[string]*TrieNode)
			}
			node.Children[part] = &TrieNode{Children: make(map[string]*TrieNode)}
		}
		node = node.Children[part]
	}
	node.IsEndOfWord = true
}

func (t *HybridTrie) RemovePath(path string) error {
	path = strings.ReplaceAll(path, "\\", "/") // Normalize path
	parts := strings.Split(path, "/")          // Split path into parts
	node := &t.Root
	for _, part := range parts {
		if _, ok := node.Children[part]; !ok {
			return errors.New("path not found")
		}
		node = node.Children[part]
	}
	if !node.IsEndOfWord {
		return errors.New("path not found")
	}
	node.IsEndOfWord = false
	return nil
}

func (t *HybridTrie) Search(filename string) []string {
	results := []string{}
	t.searchHelper(&t.Root, "", filename, &results)
	return results
}

func (t *HybridTrie) searchHelper(node *TrieNode, currentPath, filename string, results *[]string) {
	for part, child := range node.Children {
		newPath := currentPath
		if newPath != "" {
			newPath += "/"
		}
		newPath += part
		if child.IsEndOfWord && part == filename {
			*results = append(*results, newPath)
		}
		t.searchHelper(child, newPath, filename, results)
	}
}

type Match struct {
	Path    string
	Indexes []int
	Score   int
}

func (t *HybridTrie) FuzzySearch(filename string) []Match {
	paths := t.getAllPaths(&t.Root, "")
	matches := fuzzy.Find(filename, paths)
	results := make([]Match, len(matches))
	for i, match := range matches {
		results[i] = Match{
			Path:    match.Str,
			Indexes: match.MatchedIndexes,
			Score:   match.Score,
		}
	}
	return results
}

func (t *HybridTrie) getAllPaths(node *TrieNode, currentPath string) []string {
	paths := []string{}
	if node.IsEndOfWord {
		paths = append(paths, currentPath)
	}
	for part, child := range node.Children {
		newPath := currentPath
		if newPath != "" {
			newPath += "/"
		}
		newPath += part
		paths = append(paths, t.getAllPaths(child, newPath)...)
	}
	return paths
}

const MaxFilesToPrint = 25

func (t *HybridTrie) PrintTrie() {
	t.printHelper(&t.Root, "")
}

func (t *HybridTrie) printHelper(node *TrieNode, indent string) {
	if len(node.Children) > MaxFilesToPrint {
		log.Println(indent + "Large directory ")
		return
	}
	for part, child := range node.Children {
		log.Println(indent + part)
		if child.IsEndOfWord {
			log.Println(indent + "└── *")
		}
		t.printHelper(child, indent+"    ")
	}
}

func (t *HybridTrie) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	gw := gzip.NewWriter(file)
	defer gw.Close()

	encoder := gob.NewEncoder(gw)
	err = encoder.Encode(t)
	if err != nil {
		return err
	}
	return nil
}

func (t *HybridTrie) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	gr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gr.Close()

	decoder := gob.NewDecoder(gr)
	err = decoder.Decode(t)
	if err != nil {
		return err
	}
	return nil
}
