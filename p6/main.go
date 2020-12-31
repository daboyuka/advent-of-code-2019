package main

import (
	"fmt"
	"os"

	. "aoc2019/helpers"
)

type orbit struct{ Body, Sat string }

func ParseOrbits(lines []string) (orbits []orbit) {
	for _, line := range lines {
		center, orbiter := Split2(line, ")")
		orbits = append(orbits, orbit{center, orbiter})
	}
	return orbits
}

type orbitTreeNode struct {
	Body     string
	Orbits   *orbitTreeNode
	Orbiters []*orbitTreeNode
}

func ToOrbitTree(orbits []orbit) (root *orbitTreeNode, bodyToNode map[string]*orbitTreeNode) {
	bodyToNode = make(map[string]*orbitTreeNode, len(orbits))
	for _, orb := range orbits {
		n := &orbitTreeNode{Body: orb.Body}
		bodyToNode[orb.Body] = n
		if orb.Body == "COM" {
			root = n
		}
	}
	for _, orb := range orbits {
		if _, ok := bodyToNode[orb.Sat]; !ok {
			bodyToNode[orb.Sat] = &orbitTreeNode{Body: orb.Sat}
		}
	}

	for _, orb := range orbits {
		body, sat := bodyToNode[orb.Body], bodyToNode[orb.Sat]
		body.Orbiters = append(body.Orbiters, sat)
		sat.Orbits = body
	}

	return root, bodyToNode
}

func A(orbits []orbit) {
	root, _ := ToOrbitTree(orbits)

	totalDepth := 0
	var dft func(n *orbitTreeNode, depth int)
	dft = func(n *orbitTreeNode, depth int) {
		totalDepth += depth
		for _, orbiter := range n.Orbiters {
			dft(orbiter, depth+1)
		}
	}
	dft(root, 0)

	fmt.Println(totalDepth)
}

func B(orbits []orbit) {
	_, bodyToNode := ToOrbitTree(orbits)

	youBody, sanBody := bodyToNode["YOU"].Orbits, bodyToNode["SAN"].Orbits

	sanAncestors := make(map[*orbitTreeNode]int) // node -> transfers up
	for n, i := sanBody, 0; n != nil; n, i = n.Orbits, i+1 {
		sanAncestors[n] = i
	}

	for n, i := youBody, 0; n != nil; n, i = n.Orbits, i+1 {
		if sanTransfers, ok := sanAncestors[n]; ok {
			fmt.Println(sanTransfers + i)
			return
		}
	}
}

func main() {
	switch arg := os.Args[1]; arg {
	case "a":
		A(ParseOrbits(Readlines(os.Stdin)))
	case "b":
		B(ParseOrbits(Readlines(os.Stdin)))
	default:
		panic(arg)
	}
}
