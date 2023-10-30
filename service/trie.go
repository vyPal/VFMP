package main

import (
	"encoding/gob"
	"errors"
	"os"
)

type TrieNode struct {
	IsEndOfWord bool
	Children    map[string]*TrieNode
	FullPath    string
}

type HybridTrie struct {
	Root TrieNode
}

func (t *HybridTrie) AddPath(path string) {
	node := &t.Root
	for _, char := range path {
		key := string(char)
		if _, ok := node.Children[key]; !ok {
			if node.Children == nil {
				node.Children = make(map[string]*TrieNode)
			}
			node.Children[key] = &TrieNode{Children: make(map[string]*TrieNode)}
		}
		node = node.Children[key]
	}
	node.IsEndOfWord = true
	node.FullPath = path
}

func (t *HybridTrie) RemovePath(path string) error {
	node := &t.Root
	for _, char := range path {
		key := string(char)
		if _, ok := node.Children[key]; !ok {
			return errors.New("path not found")
		}
		node = node.Children[key]
	}
	if !node.IsEndOfWord {
		return errors.New("path not found")
	}
	node.IsEndOfWord = false
	node.FullPath = ""
	return nil
}

func (t *HybridTrie) ExactMatch(path string) bool {
	node := &t.Root
	for _, char := range path {
		key := string(char)
		if _, ok := node.Children[key]; !ok {
			return false
		}
		node = node.Children[key]
	}
	return node.IsEndOfWord
}

func (t *HybridTrie) FuzzyMatch(path string) []string {
	node := &t.Root
	for _, char := range path {
		key := string(char)
		if _, ok := node.Children[key]; !ok {
			return []string{}
		}
		node = node.Children[key]
	}
	result := []string{}
	if node.IsEndOfWord {
		result = append(result, node.FullPath)
	}
	for _, child := range node.Children {
		result = append(result, t.FuzzyMatchHelper(child)...)
	}
	return result
}

func (t *HybridTrie) FuzzyMatchHelper(node *TrieNode) []string {
	result := []string{}
	if node.IsEndOfWord {
		result = append(result, node.FullPath)
	}
	for _, child := range node.Children {
		result = append(result, t.FuzzyMatchHelper(child)...)
	}
	return result
}

func (t *HybridTrie) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
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
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(t)
	if err != nil {
		return err
	}
	return nil
}
