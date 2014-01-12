package LevenshteinTrie

import (
	"unicode/utf8"
)

func Min(a ...int) int {
	min := int(^uint(0) >> 1) // largest int
	for _, i := range a {
		if i < min {
			min = i
		}
	}
	return min
}
func Max(a ...int) int {
	max := int(0)
	for _, i := range a {
		if i > max {
			max = i
		}
	}
	return max
}

type TrieNode struct {
	letter rune //Equivalent to int32
	children map[rune]*TrieNode
	isFinal bool
	text string
}

func (root *TrieNode) InsertText(text string) {

	if (root == nil) {
		root = &TrieNode{children: make(map[rune]*TrieNode)}
	}

	currNode := root //Starts at root
	for i, w := 0, 0; i < len(text); i += w {
		runeValue, width := utf8.DecodeRuneInString(text[i:])
		final := false
		if (width + i == len(text)) {
			final = true
		}
		w = width

		currNode = buildTrieNode(currNode, runeValue, final, text)
	}
}

func buildTrieNode(t *TrieNode, runeValue rune, final bool, text string) (*TrieNode) {

	if node, exists := t.children[runeValue]; exists {
		if (final) {
			node.text = text
		}
		return node
	} else {
		node := &TrieNode{letter: runeValue, isFinal: final, children: make(map[rune]*TrieNode) }
		t.children[runeValue] = node
		return node
	}
	return nil
}

func (t *TrieNode) SearchSuffix(query string) ([]string) {

	var curr *TrieNode
	var ok bool
	//first, find the end of the prefix
	for _, rune := range query {
		if curr, ok = curr.children[rune]; ok {
			//do nothing
		}
		return nil
	}

	candidates := make([]string, 0)

	candidates = getAllSuffixes(curr, candidates)

	return candidates
}

func getAllSuffixes(n *TrieNode, candidates []string) ([]string){

	if (n.isFinal) {
		candidates = append(candidates, n.text)
	}

	for _, childNode := range n.children {
		candidates = getAllSuffixes(childNode, candidates)
	}

	return candidates
}

func SearchLevenshten(n *TrieNode, text string, distance int) ([]string) {
	candidates := make([]string, 0)

	//initialize the first row
	currentRow := make([]int, len(text) + 1)

	for i := 0; i < len(currentRow); i++ {
		currentRow[i] = i
	}

	for letter, childNode := range n.children {
		candidates = searchRecursive(childNode, currentRow, letter, []rune(text), distance, candidates)
	}

	return candidates
}

func searchRecursive(n *TrieNode, prevRow []int, letter rune, text []rune, maxDistance int, candidates []string) ([]string) {
    columns := len(text) + 1
    currentRow := make([]int, columns)

	currentRow[0] = prevRow[0] + 1

	for col := 1; col <  columns; col++ {
		insertCost := currentRow[col - 1] + 1
		deleteCost := currentRow[col] + 1
		var replaceCost int
		if text[col - 1] != letter {
			  if text[col - 1] != letter {
				  replaceCost = prevRow[col - 1] + 1
			  } else {
				  replaceCost = prevRow[col - 1]
			  }
		}
		currentRow[col] = Min(insertCost, deleteCost, replaceCost)
	}

	if currentRow[len(currentRow) - 1] <= maxDistance && len(n.text) > 0 {
		candidates = append(candidates, n.text)
	}

	if Min(currentRow...) <= maxDistance {
		for letter, childNode := range n.children {
			candidates = searchRecursive(childNode, currentRow, letter, []rune(text), maxDistance, candidates)
		}
	}

	return candidates
}
