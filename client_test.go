package opal

import (
	"testing"
)

func TestEverything(t *testing.T) {
	if !testing.Verbose() {
		t.Skip("-v not passed; not running TestEverything")
	}

	c, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer func() {
		if err := c.WriteConfig(); err != nil {
			t.Errorf("c.WriteConfig: %v", err)
		}
	}()

	o, err := c.Overview()
	if err != nil {
		t.Fatalf("c.Overview: %v", err)
	}
	for _, card := range o.Cards {
		t.Logf("Card %s: $%.2f", card.Number, float64(card.Balance)/100)
	}
	if len(o.Cards) == 0 {
		return
	}

	a, err := c.Activity(0)
	if err != nil {
		t.Fatalf("c.Activity(0): %v", err)
	}
	if n1, n2 := a.CardNumber, o.Cards[0].Number; n1 != n2 {
		t.Errorf("Card number from c.Activity = %q, different from c.Overview = %q", n1, n2)
	}
	for _, tr := range a.Transactions {
		t.Logf("%v\t(%5s) %s [$%.2f]", tr.When, tr.Mode, tr.Details, float64(-tr.Amount)/100)
	}
}
