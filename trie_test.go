package LevenshteinTrie

import (
	"bufio"
	"os"
	"testing"
)

var tree *TrieNode

func TestInsert(t *testing.T) {
	tree = NewTrie()
	file, err := os.Open("/usr/share/dict/words")
	if err != nil {
		t.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if err != nil {
			break
		}

		tree.InsertText(line)
	}
}

func TestPrefixSearch(t *testing.T) {
	expected := []string{
		"zygosaccharomyces",
		"zygose",
		"zygosis",
		"zygosperm",
		"zygosphenal",
		"zygosphene",
		"zygosphere",
		"zygosporange",
		"zygosporangium",
		"zygospore",
		"zygosporic",
		"zygosporophore",
		"zygostyle",
	}
	words := tree.Suffix("zygos")
	for _, e := range expected {
		if !contains(words, e) {
			t.Errorf("Missing word: %s", e)
		}
	}
}

func contains(words []string, word string) bool {
	for _, w := range words {
		if w == word {
			return true
		}
	}
	return false
}
func TestLevenshteinSearch(t *testing.T) {
	expected := []struct {
		query  string
		result []QueryResult
	}{
		{"accidia", []QueryResult{
			{"accidia", 0},
			{"accidie", 1},
		}},
	}
	for _, e := range expected {
		results := tree.Levenshtein(e.query, 1)
		if !containsq(results, e.query) {
			t.Errorf("Missing term: %s", e.query)
		}
	}
}

func containsq(results []QueryResult, word string) bool {
	for _, r := range results {
		if r.Val == word {
			return true
		}
	}
	return false
}
