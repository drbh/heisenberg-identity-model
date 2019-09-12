package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/davecgh/go-spew/spew"
)

type Assertion struct {
	UUID    int    `json:"uuid"`
	Type    string `json:"type"`
	Claimer string `json:"claimer"`
	Target  string `json:"target"`
	Psuedo  string `json:"psuedo"`
	Reason  int    `json:"reason"`
}

type Cloud struct {
	Assertions []Assertion `json:"assertions"`
}

type Link struct {
	UUID   int    `json:"uuid"`
	From   int    `json:"from"`
	FromOn string `json:"fromOn"`
	To     int    `json:"to"`
	ToOn   string `json:"toOn"`
	Type   string `json:"type"`
}

func getPsuedos(users []Assertion, x string) []Assertion {
	var startVal []Assertion
	for _, val := range users {
		if val.Psuedo == x {
			startVal = append(startVal, val)
		}
	}
	return startVal
}

func getTargets(users []Assertion, x string) []Assertion {
	var startVal []Assertion
	for _, val := range users {
		if val.Target == x {
			startVal = append(startVal, val)
		}
	}
	return startVal
}
func getUUID(users []Assertion, start int) Assertion {
	var startVal []Assertion
	for _, val := range users {
		if val.UUID == start {
			startVal = append(startVal, val)
			break
		}
	}
	return startVal[0]
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func expandGraph(start int) ([]Assertion, []Link) {
	jsonFile, err := os.Open(
		"/Users/drbh2/Desktop/heisenberg-identity-model/cloud-bootstrap.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var users []Assertion
	json.Unmarshal(byteValue, &users)

	var startVal []Assertion
	for _, val := range users {
		if val.UUID == start {
			startVal = append(startVal, val)
			break
		}
	}

	var seen []string
	var seenAssertion []Assertion
	var toexplore = make(map[string]int)

	toexplore[startVal[0].Psuedo] = 1

	for {
		if len(toexplore) == 0 {
			break
		}
		keys := make([]string, 0, len(toexplore))
		for k := range toexplore {
			keys = append(keys, k)
		}
		// only use the first key for look up
		x, keepkeys := keys[0], keys[1:]
		wpsu := getPsuedos(users, x)
		tpsu := getTargets(users, x)
		wpsu = append(wpsu, tpsu...)
		seen = append(seen, x)
		// remove that first key from out to explore
		for k := range toexplore {
			if contains(keepkeys, k) {
			} else {
				delete(toexplore, k)
			}
		}
		for _, val := range wpsu {
			if !contains(seen, val.Psuedo) {
				toexplore[val.Psuedo] = 1
				seenAssertion = append(seenAssertion, val)
			}
			if !contains(seen, val.Target) {
				toexplore[val.Target] = 1
				seenAssertion = append(seenAssertion, val)
			}
		}
	}
	// spew.Dump(seenAssertion)
	var links []Link
	counter := 1
	for _, val := range seenAssertion {
		if val.Reason == 0 {
			continue
		}
		var z Assertion
		for _, z = range seenAssertion {
			if z.UUID == val.Reason {
				break
			}
		}
		l := Link{
			counter,
			val.UUID,
			"UUID",
			z.UUID,
			"reason",
			"because-of",
		}
		counter = counter + 1
		links = append(links, l)
	}

	for _, val := range seenAssertion {
		psu := getPsuedos(seenAssertion, val.Psuedo)
		for _, p := range psu {
			if p.UUID == val.UUID || p.Reason != 0 {
				continue
			}
			// spew.Dump(p)
			l := Link{
				counter,
				val.UUID,
				"Psuedo",
				p.UUID,
				"Psuedo",
				"also-is",
			}
			counter = counter + 1
			links = append(links, l)
		}

	}

	for _, val := range seenAssertion {
		psu := getTargets(seenAssertion, val.Target)
		for _, p := range psu {
			if p.UUID == val.UUID || p.Reason == 0 {
				continue
			}
			// spew.Dump(p)
			l := Link{
				counter,
				val.UUID,
				"Target",
				p.UUID,
				"Target",
				"recognized",
			}
			counter = counter + 1
			links = append(links, l)
		}
	}
	return seenAssertion, links
}

func traverseItAll(n int) {

	seenAssertion, links := expandGraph(n)
	// spew.Dump(seenAssertion)

	assert := getUUID(seenAssertion, n)
	fmt.Println(assert.Target)

	for _, l := range links {

		pa := []Link{}
		if l.To == assert.UUID && l.Type == "because-of" {
			pa = append(pa, l)
			// spew.Dump(l)
			// assert = getUUID(seenAssertion, l.From)
			// spew.Dump(assert)
			for _, lx := range links {
				if lx.From == l.From && lx.UUID != l.UUID && lx.Type == "recognized" {
					pa = append(pa, lx)
					// spew.Dump(lx)

					for _, ly := range links {
						if ly.From == lx.To && ly.UUID != l.UUID && ly.Type == "because-of" {
							pa = append(pa, ly)
							// spew.Dump(ly)
						}
					}

				}

			}
			// break
			spew.Dump(pa)
			fmt.Println("----")
		}

	}
}

func main() {
	// c := Cloud{}
	seenAssertion, _ := expandGraph(77)
	spew.Dump(seenAssertion)

	traverseItAll(77)
	// traverseItAll(2)
	// traverseItAll(4)
	// n := 4

	// 100: D. Holtz -(dcc myid)-> 56

	// 100: D. Holtz -(dcc oauth)-> 56 -(dcc oauth)-> drbh@gmail.com
	// 090: D. Holtz -(dcc oauth)-> 56 -(dcc oauth)-> drbh@gmail.com -(dsu oauth)-> 000231-learner
	// 090: D. Holtz -(dcc oauth)-> 56 -(dcc oauth)-> drbh@gmail.com -(dsu oauth)-> 000231-learner -(dsu oauth)-> David Richard Blyn Holtz

	// 090: D. Holtz -(dcc oauth)-> 56 -(dcc oauth)-> drbh@gmail.com -(dsu oauth)-> 000231-learner -(dsu oauth)-> uport:did-000
	// 070: D. Holtz -(dcc oauth)-> 56 -(dcc oauth)-> drbh@gmail.com -(dsu oauth)-> 000231-learner -(dsu oauth)-> uport:did-000 -(asu oauth)-> student-003
	// 070: D. Holtz -(dcc oauth)-> 56 -(dcc oauth)-> drbh@gmail.com -(dsu oauth)-> 000231-learner -(dsu oauth)-> uport:did-000 -(asu oauth)-> student-003 -(asu oauth)-> David Richard Holtz

}
