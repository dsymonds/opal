package opal

import (
	"testing"
)

func TestEverything(t *testing.T) {
	if !testing.Verbose() {
		t.Skip("-v not passed; not running TestEverything")
	}

	c, err := NewClient(FileAuthStore(DefaultAuthFile))
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
		t.Logf("Card \"%s\" balance: $%.2f", card.Name, float64(card.Balance)/100)
	}
	if len(o.Cards) == 0 {
		return
	}

	a, err := c.Activity(0)
	if err != nil {
		t.Fatalf("c.Activity(0): %v", err)
	}
	if n1, n2 := a.CardName, o.Cards[0].Name; n1 != n2 {
		t.Errorf("Card name from c.Activity = %q, different from c.Overview = %q", n1, n2)
	}
	for _, tr := range a.Transactions {
		t.Logf("%d %v\t(%5s) %s [$%.2f] [$%.2f] [$%.2f]", tr.JourneyNumber, tr.When, tr.Mode, tr.Details, float64(tr.Fare)/100, float64(tr.Discount)/100, float64(-tr.Amount)/100)
	}
}
