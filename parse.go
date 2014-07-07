package opal

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
)

type Overview struct {
	Cards []Card
}

type Card struct {
	Number  string
	Balance int // in cents
}

var (
	cardNumberRE = regexp.MustCompile(`^\d+$`)
	balanceRE    = regexp.MustCompile(`^\$(\d+)\.(\d\d)$`)
)

// parseCard parses card info from the num and bal TDs.
func parseCard(num, bal string) (Card, error) {
	if !cardNumberRE.MatchString(num) {
		return Card{}, fmt.Errorf("card number %q does not match /%v/", num, cardNumberRE)
	}
	m := balanceRE.FindStringSubmatch(bal)
	if m == nil {
		return Card{}, fmt.Errorf("balance %q does not match /%v/", bal, balanceRE)
	}
	dollars, err := strconv.ParseInt(m[1], 10, 32)
	if err != nil {
		return Card{}, fmt.Errorf("bad balance %q: %v", bal, err)
	}
	cents, _ := strconv.ParseInt(m[2], 10, 8) // can't fail; it's exactly two digits
	return Card{
		Number:  num,
		Balance: int(dollars)*100 + int(cents),
	}, nil
}

// parseOverview parses a page fetched from https://www.opal.com.au/registered/index.
func parseOverview(input []byte) (*Overview, error) {
	// TODO: check input is UTF-8
	doc, err := html.Parse(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}

	activeTable := findByAttr(doc, "id", "dashboard-active-cards")
	if activeTable == nil || activeTable.DataAtom != atom.Table {
		return nil, errors.New("did not find active table")
	}
	tbody := findByDataAtom(activeTable, atom.Tbody)
	if tbody == nil {
		return nil, errors.New("did not find tbody")
	}

	var cardRows [][]string // one per row, each row having two elements (number and balance)
	eachByAtom(tbody, atom.Tr, func(n *html.Node) bool {
		var tds []string
		eachByAtom(n, atom.Td, func(n *html.Node) bool {
			if len(tds) < 2 {
				tds = append(tds, text(n))
			}
			return false
		})
		if len(tds) == 2 {
			cardRows = append(cardRows, tds)
			return false
		}
		return false
	})

	o := new(Overview)
	for _, row := range cardRows {
		card, err := parseCard(row[0], row[1])
		if err != nil {
			return nil, fmt.Errorf("parsing card row: %v", err)
		}
		o.Cards = append(o.Cards, card)
	}
	return o, nil
}

func parseLogin(input []byte) (token string, err error) {
	// TODO: check input is UTF-8
	doc, err := html.Parse(bytes.NewReader(input))
	if err != nil {
		return "", err
	}
	node := findByAttr(doc, "name", "CSRFToken")
	if node == nil {
		return "", errors.New("did not find CSRFToken <input>")
	}
	for _, attr := range node.Attr {
		if attr.Key == "value" {
			return attr.Val, nil
		}
	}
	return "", fmt.Errorf("unexpected form of CSRFToken: %s", render(node))
}

func text(n *html.Node) string {
	n = find(n, func(n *html.Node) bool { return n.Type == html.TextNode })
	if n == nil {
		return ""
	}
	return n.Data
}

func findByAttr(n *html.Node, key, val string) *html.Node {
	return find(n, func(n *html.Node) bool {
		for _, attr := range n.Attr {
			if attr.Key == key && attr.Val == val {
				return true
			}
		}
		return false
	})
}

func findByDataAtom(n *html.Node, a atom.Atom) *html.Node {
	return find(n, func(n *html.Node) bool {
		return n.DataAtom == a
	})
}

func eachByAtom(n *html.Node, a atom.Atom, f func(*html.Node) bool) {
	each(n, func(n *html.Node) bool {
		if n.DataAtom == a {
			f(n)
		}
		return true
	})
}

// each calls f for each node from n downwards in a BFS order.
// If f returns false then that node's children will be skipped.
func each(n *html.Node, f func(*html.Node) bool) {
	if !f(n) {
		return
	}
	for kid := n.FirstChild; kid != nil; kid = kid.NextSibling {
		each(kid, f)
	}
}

func find(n *html.Node, pred func(*html.Node) bool) *html.Node {
	if pred(n) {
		return n
	}
	for kid := n.FirstChild; kid != nil; kid = kid.NextSibling {
		if n := find(kid, pred); n != nil {
			return n
		}
	}
	return nil
}

func render(n *html.Node) string {
	var buf bytes.Buffer
	html.Render(&buf, n)
	return buf.String()
}
