package LevenshteinTrie

import (
	"bufio"
	"errors"
	"io"
	"os"
	"testing"
)

var tree *TrieNode

func getfile() (io.ReadCloser, error) {
	filename := "./w1_fixed.txt"
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		file, err := os.Open(filename)
		return file, err
	} else {
		file, err := os.Open("/usr/share/dict/words")
		return file, err
	}
	return nil, errors.New("Dictionary file does not exist")
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
		query    string
		result   []QueryResult
		distance int
	}{
		{"accidia", []QueryResult{
			{"accidia", 0},
			{"accidie", 1},
		}, 1},
	}
	for _, e := range expected {
		results := tree.Levenshtein(e.query, e.distance)
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
