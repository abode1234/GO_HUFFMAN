package main

import (
	"container/heap"
	"fmt"
	"io/ioutil"
	"log"
)

// HuffmanNode represents a node in the Huffman Tree
type HuffmanNode struct {
	character rune
	frequency int
	left      *HuffmanNode
	right     *HuffmanNode
}

// HuffmanHeap implements heap.Interface for HuffmanNode
type HuffmanHeap []*HuffmanNode

func (h HuffmanHeap) Len() int           { return len(h) }
func (h HuffmanHeap) Less(i, j int) bool { return h[i].frequency < h[j].frequency }
func (h HuffmanHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *HuffmanHeap) Push(x interface{}) {
	*h = append(*h, x.(*HuffmanNode))
}
func (h *HuffmanHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// BuildFrequencyTable builds a frequency table for the input text
func BuildFrequencyTable(text string) map[rune]int {
	frequency := make(map[rune]int)
	for _, char := range text {
		frequency[char]++
	}
	return frequency
}

// BuildHuffmanTree builds a Huffman tree based on character frequencies
func BuildHuffmanTree(frequency map[rune]int) *HuffmanNode {
	h := &HuffmanHeap{}
	heap.Init(h)

	// Create a leaf node for each character and push it into the priority queue
	for char, freq := range frequency {
		heap.Push(h, &HuffmanNode{character: char, frequency: freq})
	}

	// Build the Huffman tree
	for h.Len() > 1 {
		left := heap.Pop(h).(*HuffmanNode)
		right := heap.Pop(h).(*HuffmanNode)
		heap.Push(h, &HuffmanNode{
			frequency: left.frequency + right.frequency,
			left:      left,
			right:     right,
		})
	}

	return heap.Pop(h).(*HuffmanNode)
}

// GenerateHuffmanCodes generates Huffman codes by traversing the tree
func GenerateHuffmanCodes(node *HuffmanNode, prefix string, codes map[rune]string) {
	if node == nil {
		return
	}
	if node.left == nil && node.right == nil {
		codes[node.character] = prefix
		return
	}
	GenerateHuffmanCodes(node.left, prefix+"0", codes)
	GenerateHuffmanCodes(node.right, prefix+"1", codes)
}

// Encode encodes the input text using Huffman codes
func Encode(text string, codes map[rune]string) string {
	encoded := ""
	for _, char := range text {
		encoded += codes[char]
	}
	return encoded
}

// Decode decodes the binary string using the Huffman tree
func Decode(encoded string, root *HuffmanNode) string {
	decoded := ""
	node := root
	for _, bit := range encoded {
		if bit == '0' {
			node = node.left
		} else {
			node = node.right
		}
		if node.left == nil && node.right == nil {
			decoded += string(node.character)
			node = root
		}
	}
	return decoded
}

// ReadFile reads the content of a file and returns it as a string
func ReadFile(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	return string(content)
}

// WriteToFile writes a string to a file
func WriteToFile(filename, content string) {
	err := ioutil.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}
}

func main() {
	// Read the input text from input.txt
	inputText := ReadFile("input.txt")

	// Build frequency table and Huffman Tree
	frequency := BuildFrequencyTable(inputText)
	huffmanTree := BuildHuffmanTree(frequency)

	// Generate Huffman Codes
	codes := make(map[rune]string)
	GenerateHuffmanCodes(huffmanTree, "", codes)

	// Encode the input text
	encoded := Encode(inputText, codes)
	// Write the encoded text to encoded.txt
	WriteToFile("encoded.txt", encoded)

	// Decode the encoded text
	decoded := Decode(encoded, huffmanTree)
	// Write the decoded text to decoded.txt
	WriteToFile("decoded.txt", decoded)

	fmt.Println("Encoding and decoding complete. Check encoded.txt and decoded.txt files.")
}

