package LevenshteinTrie

import (
	"bufio"
	"io"
	"os"
	"testing"
)

var tree *TrieNode

func getfile() (io.ReadCloser, error) {
	filename := "./w1_fixed.txt"
	file, err := os.Open(filename)
	return file, err
}

func TestInsert(t *testing.T) {
	tree = NewTrie()
	file, err := getfile()
	defer file.Close()
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
		"zygodactyl",
		"zygoma",
		"zygomatic",
		"zygomaticus",
		"zygomycetes",
		"zygon",
		"zygosity",
		"zygote",
		"zygote-specific",
		"zygotes",
		"zygourakis",
	}
	query := "zygo"
	words := tree.Suffix(query)
	for _, e := range expected {
		if !contains(words, e) {
			t.Errorf("Missing word: %s", e)
		} else {
			t.Logf("Found: %s\n", e)
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
		query    string
		result   []QueryResult
		distance int
	}{
		{"accidens", []QueryResult{
			{"accidens", 0},
			{"accident", 1},
		}, 1},
	}
	for _, e := range expected {
		results := tree.Levenshtein(e.query, e.distance)
		if !containsq(results, e.query) {
			t.Errorf("Missing term: %s", e.query)
		} else {
			t.Logf("Looking for: %s, Found: %s\n", e.query, results)
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
