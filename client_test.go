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
}
