package opal

import (
	"reflect"
	"testing"
	"time"
)

func TestParseAmount(t *testing.T) {
	tests := []struct {
		in   string
		want int
	}{
		{"$0.00", 0},
		{"$100.00", 10000},
		{"$4.10", 410},
		{"-$4.10", -410},
	}
	for _, tc := range tests {
		got, err := parseAmount(tc.in)
		if err != nil {
			t.Errorf("parseAmount(%q): %v", tc.in, err)
			continue
		}
		if got != tc.want {
			t.Errorf("parseAmount(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestParseLogin(t *testing.T) {
	token, err := parseLogin([]byte(loginPage))
	if err != nil {
		t.Fatalf("parseLogin: %v", err)
	}
	const want = "xxx-yyy-zzz"
	if token != want {
		t.Errorf("parsed login token: got %q, want %q", token, want)
	}
}

const loginPage = `
<input value="Log in" type="submit" tabindex="21"></span></div><a title="Forgot your username or password?" href="/login/forgotten" tabindex="22">Forgot your username or password?</a></fieldset><div><input type="hidden" name="CSRFToken" value="xxx-yyy-zzz" tabindex="-1">
`

func TestParseOverview(t *testing.T) {
	o, err := parseOverview([]byte(overviewPage))
	if err != nil {
		t.Fatalf("parseOverview: %v", err)
	}
	want := &Overview{
		Cards: []Card{{
			Number:  "31415926535",
			Balance: 7743,
		}},
	}
	if !reflect.DeepEqual(o, want) {
		t.Errorf("parseOverview returned incorrect data.\n got %+v\nwant %+v", o, want)
	}
}

const overviewPage = `<html>
<table id="dashboard-active-cards"><caption><span>My Opal cards</span></caption><thead><tr><th>Opal Card</th><th>Balance</th><th>Status</th><th></th></tr></thead><tbody><tr><td id="nameCol0">31415926535</td><td>$77.43</td><td id="activateCol0"><a title="Activate Opal card" href="0" class="markActivation">Activate Opal card</a></td><td><a title="View" class="dashboard-card-view" href="index?cardIndex=0">View</a></td></tr></tbody></table>
`

func TestParseActivity(t *testing.T) {
	a, err := parseActivity([]byte(activityPage))
	if err != nil {
		t.Fatalf("parseActivity: %v", err)
	}
	want := &Activity{
		CardNumber: "31415926535",
		Transactions: []*Transaction{
			{
				Number:      5,
				When:        time.Date(2014, time.July, 9, 17, 1, 0, 0, sydneyZone),
				Mode:        "train",
				Details:     "Town Hall to No tap off",
				FareApplied: "Default fare",
				Fare:        810,
				Amount:      -810,
			},
			{
				Number:        3,
				When:          time.Date(2014, time.July, 9, 7, 49, 0, 0, sydneyZone),
				Mode:          "train",
				Details:       "Chatswood to Town Hall",
				JourneyNumber: 1,
				Fare:          410,
				Amount:        -410,
			},
			{
				Number:  2,
				When:    time.Date(2014, time.July, 9, 7, 49, 0, 0, sydneyZone),
				Details: "Top up - opal.com.au",
				Amount:  10000,
			},
		},
	}
	if !reflect.DeepEqual(a, want) {
		t.Errorf("parseActivity returned incorrect data.\n got %+v\nwant %+v", a, want)
	}
}

const activityPage = `<html>
<table id="transaction-data"><caption><span>My Opal activity: 31415926535</span></caption>
<thead><tr><th>Transaction<br/>number</th><th>Date/time</th><th class="narrow center">Mode</th><th>Details</th><th class="narrow center">Journey<br/>number</th><th class="right">Fare Applied</th><th class="right">Fare</th><th class="right amount">Discount</th><th class="right amount">Amount</th></tr></thead>
<tbody style="opacity: 1;">
<tr><td>5</td><td class="date-time">Wed<br/>09/07/2014<br/>17:01</td><td class="center"><img height="32" width="32" alt="train" src="/images/icons/mode-train.png"/></td><td class="transaction-summary">Town Hall to No tap off </td><td></td><td class="right">Default fare</td><td class="right nowrap">$8.10</td><td class="right nowrap">$0.00</td><td class="right nowrap">-$8.10</td></tr>
<tr><td>3</td><td class="date-time">Wed<br/>09/07/2014<br/>07:49</td><td class="center"><img height="32" width="32" alt="train" src="/images/icons/mode-train.png"/></td><td class="transaction-summary">Chatswood to Town Hall</td><td>1</td><td class="right"></td><td class="right nowrap">$4.10</td><td class="right nowrap">$0.00</td><td class="right nowrap">-$4.10</td></tr>
<tr><td>2</td><td class="date-time">Wed<br/>09/07/2014<br/>07:49</td><td class="center"></td><td class="transaction-summary">Top up - opal.com.au</td><td></td><td class="right"></td><td class="right nowrap"></td><td class="right nowrap"></td><td class="right nowrap">$100.00</td></tr>
</tbody></table>
`
