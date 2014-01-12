LevenshteinTrie
===============

A Trie data structure that allows for fuzzy string matching

This is the Go version of a python program written by Steve Hanov in his [blog post](http://stevehanov.ca/blog/index.php?id=114)

It is not finished or tested yet.

###How it works

 - It is a basic [Trie](http://en.wikipedia.org/wiki/Trie).

 - You can search for all words that are suffixes of a string. 

 - You can also search for words within a certain edit distance of a string. The algorithm memoizes the Levenshtein algorithm when it recursively iterates through the Trie nodes. This speeds up the Levenshtein matches hugely.
