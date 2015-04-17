package LevenshteinTrie

import (
	"fmt"
	"sort"
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
	letter   rune //Equivalent to int32
	children map[rune]*TrieNode
	final    bool
	text     string
}

func NewTrie() *TrieNode {
	return &TrieNode{children: make(map[rune]*TrieNode)}
}

func (root *TrieNode) InsertText(text string) {

	if root == nil {
		return
	}

	currNode := root //Starts at root
	for i, w := 0, 0; i < len(text); i += w {
		runeValue, width := utf8.DecodeRuneInString(text[i:])
		final := false
		if width+i == len(text) {
			final = true
		}
		w = width

		currNode = NewTrieNode(currNode, runeValue, final, text)
	}
}

func NewTrieNode(t *TrieNode, runeValue rune, final bool, text string) *TrieNode {
	node, exists := t.children[runeValue]
	if exists {
		if final {
			node.final = true
			node.text = text
		}
		return node
	} else {
		node = &TrieNode{letter: runeValue, children: make(map[rune]*TrieNode)}
		t.children[runeValue] = node
		return node
	}
	return nil
}

func (t *TrieNode) SearchSuffix(query string) []string {

	var curr *TrieNode
	var ok bool
	//first, find the end of the prefix
	for _, letter := range query {
		if curr != nil {
			if curr, ok = curr.children[letter]; ok {
				//do nothing
			}

		}
	}

	candidates := make([]string, 0)

	var getAllSuffixes func(n *TrieNode)
	getAllSuffixes = func(n *TrieNode) {
		if n == nil {
			return
		}
		if n.final == true {
			candidates = append(candidates, n.text)
		}

		for _, childNode := range n.children {
			getAllSuffixes(childNode)
		}

	}
	getAllSuffixes(curr)

	return candidates
}

type QueryResult struct {
	Val      string
	Distance int
}

func (q QueryResult) String() string {
	return fmt.Sprintf("Val: %s\n", q.Val)
}

type ByDistance []QueryResult

func (a ByDistance) Len() int           { return len(a) }
func (a ByDistance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDistance) Less(i, j int) bool { return a[i].Distance < a[j].Distance }

func (n *TrieNode) SearchLevenshtein(text string, distance int) []QueryResult {

	//initialize the first row for the dynamic programming alg
	l := utf8.RuneCount([]byte(text))
	currentRow := make([]int, l+1)

	for i := 0; i < len(currentRow); i++ {
		currentRow[i] = i
	}

	candidates := make([]QueryResult, 0)

	var searchRecursive func(n *TrieNode, prevRow []int, letter rune, text []rune, maxDistance int)
	searchRecursive = func(n *TrieNode, prevRow []int, letter rune, text []rune, maxDistance int) {
		columns := len(text) + 1
		currentRow := make([]int, columns)

		currentRow[0] = prevRow[0] + 1

		for col := 1; col < columns; col++ {
			insertCost := currentRow[col-1] + 1
			deleteCost := currentRow[col] + 1
			var replaceCost int
			if text[col-1] != letter {
				if text[col-1] != letter {
					replaceCost = prevRow[col-1] + 1
				} else {
					replaceCost = prevRow[col-1]
				}
			}
			currentRow[col] = Min(insertCost, deleteCost, replaceCost)
		}

		distance := currentRow[len(currentRow)-1]
		if distance <= maxDistance && len(n.text) > 0 {
			candidates = append(candidates, QueryResult{Val: n.text, Distance: distance})
		}

		if Min(currentRow...) <= maxDistance {
			for letter, childNode := range n.children {
				searchRecursive(childNode, currentRow, letter, []rune(text), maxDistance)
			}
		}
	}

	for letter, childNode := range n.children {
		searchRecursive(childNode, currentRow, letter, []rune(text), distance)
	}
	sort.Sort(ByDistance(candidates))
	return candidates
}
