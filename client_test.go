package opal

import (
	"strings"
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
		t.Logf("Card %s: $%.2f", card.Name, float64(card.Balance)/100)
	}
	if len(o.Cards) == 0 {
		return
	}

	const lastPage = 0 // switch to 1 to see more data
	for page := 0; page <= lastPage; page++ {
		req := ActivityRequest{CardIndex: 0, Offset: page}
		a, err := c.Activity(req)
		if err != nil {
			t.Fatalf("c.Activity(%+v): %v", req, err)
		}
		if n1, n2 := a.CardName, o.Cards[0].Name; n1 != n2 {
			t.Errorf("Card name from c.Activity = %q, different from c.Overview = %q", n1, n2)
		}
		var prevWeek int
		for i, tr := range a.Transactions {
			_, week := tr.When.ISOWeek()
			if i > 0 && week != prevWeek {
				t.Logf(strings.Repeat("-", 50))
			}
			prevWeek = week
			ts := tr.When.Format("2006-01-02 15:04")
			t.Logf("%v\tJ%02d (%5s) %s [$%.2f]", ts, tr.JourneyNumber, tr.Mode, tr.Details, float64(-tr.Amount)/100)
		}
	}
}
