package opal

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
)

type Overview struct {
	Cards []Card
}

type Card struct {
	Name    string // either a name or number
	Balance int    // in cents
}

var (
	amountRE = regexp.MustCompile(`^(-?)\$(\d+)\.(\d\d)$`)
)

// parseAmount parses something matching amountRE and returns the number of cents.
func parseAmount(amt string) (int, error) {
	m := amountRE.FindStringSubmatch(amt)
	if m == nil {
		return 0, fmt.Errorf("does not match /%v/", amountRE)
	}
	dollars, err := strconv.ParseInt(m[2], 10, 32)
	if err != nil {
		return 0, err
	}
	cents, _ := strconv.ParseInt(m[3], 10, 8) // can't fail; it's exactly two digits
	x := int(dollars)*100 + int(cents)
	if m[1] == "-" {
		x = -x
	}
	return x, nil
}

// parseCard parses card info from the name and bal TDs.
func parseCard(name, bal string) (Card, error) {
	// Cards can be renamed to almost anything.
	balance, err := parseAmount(bal)
	if err != nil {
		return Card{}, fmt.Errorf("bad balance %q: %v", bal, err)
	}
	return Card{
		Name:    name,
		Balance: balance,
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

type Activity struct {
	CardName     string
	Transactions []*Transaction
}

type Transaction struct {
	Number        int
	When          time.Time
	Mode          string // "train", etc., if known
	Details       string
	JourneyNumber int // if known; numbered within the week

	FareApplied            string // e.g. "Off-peak"
	Fare, Discount, Amount int    // in cents
}

func (t *Transaction) String() string { return fmt.Sprintf("%+v", *t) }

func parseActivity(input []byte) (*Activity, error) {
	// TODO: check input is UTF-8
	doc, err := html.Parse(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}

	table := findByAttr(doc, "id", "transaction-data")
	if table == nil || table.DataAtom != atom.Table {
		return nil, errors.New("did not find transaction table")
	}

	a := new(Activity)

	// The <caption> of the table will look like
	//	<caption><span>My Opal activity: 314159</span></caption>
	caption := findByDataAtom(table, atom.Caption)
	if caption == nil || caption.FirstChild == nil {
		return nil, errors.New("did not find caption")
	}
	a.CardName = text(caption.FirstChild)
	if i := strings.Index(a.CardName, ":"); i >= 0 {
		a.CardName = strings.TrimSpace(a.CardName[i+1:])
	}

	tbody := findByDataAtom(table, atom.Tbody)
	if tbody == nil {
		return nil, errors.New("did not find tbody")
	}

	eachByAtom(tbody, atom.Tr, func(n *html.Node) bool {
		if err != nil {
			return false
		}
		var t *Transaction
		t, err = parseTransaction(n)
		a.Transactions = append(a.Transactions, t)
		return false
	})
	if err != nil {
		return nil, err
	}

	return a, nil
}

func parseTransaction(n *html.Node) (*Transaction, error) {
	// Collate all the <TD> contents.
	var tds []string
	for kid := n.FirstChild; kid != nil; kid = kid.NextSibling {
		if kid.FirstChild != nil && kid.FirstChild.DataAtom == atom.Img {
			// <td><img alt="train" ...></td>
			tds = append(tds, attrVal(kid.FirstChild, "alt"))
			continue
		}
		tds = append(tds, text(kid))
	}
	if len(tds) != 9 {
		return nil, fmt.Errorf("transaction row with %d TDs, want 9", len(tds))
	}

	// Parse the transaction.
	t := new(Transaction)
	var err error

	t.Number, err = parseDecimal(tds[0])
	if err != nil {
		return nil, fmt.Errorf("bad transaction number %q: %v", tds[0], err)
	}
	t.When, err = time.ParseInLocation("Mon 02/01/2006 15:04", tds[1], sydneyZone)
	if err != nil {
		return nil, fmt.Errorf("bad time %q: %v", tds[1], err)
	}
	t.Mode, t.Details = tds[2], strings.TrimSpace(tds[3])
	t.FareApplied = strings.TrimSpace(tds[5])

	// The rest are all optional.
	fields := []struct {
		index  int
		dst    *int
		parser func(string) (int, error)
		name   string
	}{
		{4, &t.JourneyNumber, parseDecimal, "journey number"},
		{6, &t.Fare, parseAmount, "fare"},
		{7, &t.Discount, parseAmount, "discount"},
		{8, &t.Amount, parseAmount, "amount"},
	}
	for _, f := range fields {
		if s := tds[f.index]; s != "" {
			*f.dst, err = f.parser(s)
			if err != nil {
				return nil, fmt.Errorf("bad %s %q: %v", f.name, s, err)
			}
		}
	}

	return t, nil
}

var sydneyZone *time.Location // every time is in Australia/Sydney

func init() {
	var err error
	sydneyZone, err = time.LoadLocation("Australia/Sydney")
	if err != nil {
		panic(err)
	}
}

func parseDecimal(s string) (int, error) {
	x, err := strconv.ParseInt(s, 10, 0)
	return int(x), err
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
	var bits []string
	each(n, func(n *html.Node) bool {
		if n.Type == html.TextNode {
			bits = append(bits, n.Data)
		}
		return true
	})
	return strings.Join(bits, " ")
}

func attrVal(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
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
			return f(n)
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
