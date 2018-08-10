package main

import "math/rand"

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
		for _, v := range markov {
			if v.word == str {
				return &v
			}
		}

		panic("If you are reading this message, the apocalypse has begun")
	} ("FRONT")

	for currentNode.word != "BACK" {
		out = append(out, []byte(currentNode.word)...)
		currentNode = currentNode.next()
	}

	return string(out)
}