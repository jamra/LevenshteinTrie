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

func (t *TrieNode) String() string {
	s := fmt.Sprintf("%#U\n", t.letter)
	for _, v := range t.children {
		s += fmt.Sprintf("-%#U\n", v)
	}
	return s
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
	if !exists {
		node = &TrieNode{letter: runeValue, children: make(map[rune]*TrieNode)}
		t.children[runeValue] = node
	}
	if final {
		node.final = true
		node.text = text
	}
	return node
}

func (t *TrieNode) Suffix(query string) []string {
	var curr *TrieNode
	var ok bool

	curr = t
	//first, find the end of the prefix
	for _, letter := range query {
		if curr != nil {
			curr, ok = curr.children[letter]
			if ok {
				//do nothing
			}

		} else {
			return nil
		}
	}

	candidates := getsuffixr(curr)

	return candidates
}

func getsuffixr(n *TrieNode) []string {
	if n == nil {
		return nil
	}

	candidates := make([]string, 0)
	if n.final == true {
		candidates = append(candidates, n.text)
	}

	for _, childNode := range n.children {
		candidates = append(candidates, getsuffixr(childNode)...)
	}
	return candidates
}

type QueryResult struct {
	Val      string
	Distance int
}

func (q QueryResult) String() string {
	return fmt.Sprintf("Val: %s, Dist: %d\n", q.Val, q.Distance)
}

type ByDistance []QueryResult

func (a ByDistance) Len() int           { return len(a) }
func (a ByDistance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDistance) Less(i, j int) bool { return a[i].Distance < a[j].Distance }

func (n *TrieNode) Levenshtein(text string, distance int) []QueryResult {

	//initialize the first row for the dynamic programming alg
	l := utf8.RuneCount([]byte(text))
	currentRow := make([]int, l+1)

	for i := 0; i < len(currentRow); i++ {
		currentRow[i] = i
	}

	candidates := make([]QueryResult, 0)

	for letter, childNode := range n.children {
		candidates = append(candidates, searchlevr(childNode, currentRow, letter, []rune(text), distance)...)
	}

	sort.Sort(ByDistance(candidates))
	return candidates
}

func searchlevr(n *TrieNode, prevRow []int, letter rune, text []rune, maxDistance int) []QueryResult {
	columns := len(prevRow)
	currentRow := make([]int, columns)

	currentRow[0] = prevRow[0] + 1

	for col := 1; col < columns; col++ {
		if text[col-1] == letter {
			currentRow[col] = prevRow[col-1]
			continue
		}
		insertCost := currentRow[col-1] + 1
		deleteCost := prevRow[col] + 1
		replaceCost := prevRow[col-1] + 1

		currentRow[col] = Min(insertCost, deleteCost, replaceCost)
	}

	candidates := make([]QueryResult, 0)

	distance := currentRow[len(currentRow)-1]
	if distance <= maxDistance && n.final == true {
		candidates = append(candidates, QueryResult{Val: n.text, Distance: distance})
	}
	mi := Min(currentRow[1:]...)
	if mi <= maxDistance {
		for l, childNode := range n.children {
			candidates = append(candidates, searchlevr(childNode, currentRow, l, text, maxDistance)...)
		}
	}
	return candidates
}
