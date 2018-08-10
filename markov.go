package main

import (
	"math/rand"
	"unicode"
)

type node struct {
	word string
	links []*node
}

var markov = make([]node, 0)

func (n *node) addLink(target *node) {
	for _, v := range n.links {
		if target == v {
			return
		}
	}

	n.links = append(n.links, target)
}

func (n * node) next() *node {
	return n.links[rand.Intn(len(n.links))]
}

func generate() string {
	var out []byte
	// Find FRONT
	currentNode := func(str string) *node {
		for i := range markov {
			if markov[i].word == str {
				return &markov[i]
			}
		}

		panic("If you are reading this message, the apocalypse has begun")
	} ("FRONT")
	currentNode = currentNode.next()

	for currentNode.word != "BACK" {
		out = append(out, []byte(currentNode.word)...)
		out = append(out, byte(' '))
		currentNode = currentNode.next()
	}

	// Ugly line, just makes the first character uppercase, nothing more
	out[0] = byte(unicode.ToUpper(rune(out[0])))

	// A side-effect of the above process is that there is a trailing space. This line just removes that
	out = out[:len(out)-1]

	// Add a full-stop if there isn't already one
	if out[len(out)-1] != byte('.') {
		out = append(out, byte('.'))
	}

	return string(out)
}