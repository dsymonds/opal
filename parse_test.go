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
			Name:    "My 31415926535 card",
			Balance: 7743,
		}},
	}
	if !reflect.DeepEqual(o, want) {
		t.Errorf("parseOverview returned incorrect data.\n got %+v\nwant %+v", o, want)
	}
}

const overviewPage = `<html>
<table class="dashboard-cards" id="dashboard-active-cards"><caption><span>My Opal cards</span></caption><thead><tr><th>View</th><th>Opal Card</th><th>Type</th><th>Balance</th><th>Status</th></tr></thead><tbody><tr class="alt last"><td class="bl"><input value="0" checked="checked" name="registered_card" class="card-radio-selection" id="card_0" type="radio" tabindex="43"></td><td id="nameCol0"><label for="card_0">My 31415926535 card</label></td><td>Adult</td><td>$77.43</td><td class="br">Active</td></tr></tbody></table>
`

func TestParseActivity(t *testing.T) {
	a, err := parseActivity([]byte(activityPage))
	if err != nil {
		t.Fatalf("parseActivity: %v", err)
	}
	want := &Activity{
		CardName: "31415926535 is pi",
		Transactions: []*Transaction{
			{
				Number:        6,
				When:          time.Date(2015, time.September, 29, 7, 47, 0, 0, sydneyZone),
				Mode:          "bus",
				Details:       "Willoughby Rd nr Garland to York St nr Margaret St",
				JourneyNumber: 6,
				Fare:          350,
				Amount:        -350,
			},
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
<table id="transaction-data"><caption><span>My Opal activity: 31415926535 is pi</span></caption>
<thead><tr><th>Transaction<br>number</th><th>Date/time</th><th class="narrow center">Mode</th><th>Details</th><th class="narrow center">Journey<br>number</th><th>Fare Applied</th><th class="right">Fare</th><th class="right amount">Discount</th><th class="right amount">Amount</th></tr></thead>
<tbody style="opacity: 1;">
<tr class="alt"><td>6</td><td class="date-time">Tue<br>29/09/2015<br>07:47</td><td class="center"><img height="32" width="32" alt="bus" src="/images/icons/mode-bus.png"></td><td lang="en-gb" class="transaction-summary hyphenate">Wil&shy;loughby Rd nr Gar&shy;land to York St nr Mar&shy;garet St</td><td class="center">
			    				6
										    			</td><td></td><td class="right nowrap">$3.50</td><td class="right nowrap">$0.00</td><td class="right nowrap">-$3.50</td></tr>

<tr><td>5</td><td class="date-time">Wed<br/>09/07/2014<br/>17:01</td><td class="center"><img height="32" width="32" alt="train" src="/images/icons/mode-train.png"/></td><td class="transaction-summary">Town Hall to No tap off </td><td></td><td class="right">Default fare</td><td class="right nowrap">$8.10</td><td class="right nowrap">$0.00</td><td class="right nowrap">-$8.10</td></tr>
<tr><td>3</td><td class="date-time">Wed<br/>09/07/2014<br/>07:49</td><td class="center"><img height="32" width="32" alt="train" src="/images/icons/mode-train.png"/></td><td class="transaction-summary">Chatswood to Town Hall</td><td>1</td><td class="right"></td><td class="right nowrap">$4.10</td><td class="right nowrap">$0.00</td><td class="right nowrap">-$4.10</td></tr>
<tr><td>2</td><td class="date-time">Wed<br/>09/07/2014<br/>07:49</td><td class="center"></td><td class="transaction-summary">Top up - opal.com.au</td><td></td><td class="right"></td><td class="right nowrap"></td><td class="right nowrap"></td><td class="right nowrap">$100.00</td></tr>
</tbody></table>
`
