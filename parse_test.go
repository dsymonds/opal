package opal

import (
	"reflect"
	"testing"
)

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
