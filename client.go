package opal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var configFile = filepath.Join(os.Getenv("HOME"), ".opal")

var cookieBaseURL = &url.URL{
	Scheme: "https",
	Host:   "www.opal.com.au",
}

type Client struct {
	hc *http.Client

	username, password string
}

func NewClient() (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	cfg, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	lines := bytes.Split(cfg, []byte("\n"))
	if len(lines) < 2 {
		return nil, fmt.Errorf("%s is too short; should have at least two lines", configFile)
	}
	var cookies []*http.Cookie
	for i, line := range lines[2:] {
		if len(line) == 0 {
			continue
		}
		cookie := new(http.Cookie)
		if err := json.Unmarshal(line, cookie); err != nil {
			return nil, fmt.Errorf("bad cookie at %s:%d: %v", configFile, i+3, err)
		}
		cookies = append(cookies, cookie)
	}
	jar.SetCookies(cookieBaseURL, cookies)

	c := &Client{
		hc: &http.Client{
			Jar: jar,
		},
		username: string(lines[0]),
		password: string(lines[1]),
	}
	c.hc.CheckRedirect = c.checkRedirect
	return c, nil
}

func (c *Client) WriteConfig() error {
	lines := []string{c.username, c.password}
	cookies := c.hc.Jar.Cookies(cookieBaseURL)
	for _, cookie := range cookies {
		line, err := json.Marshal(cookie)
		if err != nil {
			return err
		}
		lines = append(lines, string(line))
	}
	return ioutil.WriteFile(configFile, []byte(strings.Join(lines, "\n")+"\n"), 0600)
}

func (c *Client) Overview() (*Overview, error) {
	body, err := c.get("https://www.opal.com.au/registered/index")
	if err != nil {
		return nil, err
	}
	return parseOverview(body)
}

var errRedirect = errors.New("internal error: login redirect detected")

func (c *Client) checkRedirect(req *http.Request, via []*http.Request) error {
	if strings.HasPrefix(req.URL.Path, "/login/") {
		return errRedirect
	}
	return fmt.Errorf("hit redirect for %v", req.URL) // shouldn't happen
}

func (c *Client) get(url string) (body []byte, err error) {
	var resp *http.Response
	for try := 1; try <= 2; try++ {
		resp, err = c.hc.Get(url)
		if err != nil && strings.Contains(err.Error(), errRedirect.Error()) { // XXX: ugh, why is the value not being passed back?!?
			if err = c.login(); err == nil {
				continue // next try
			}
		}
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}
	body, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err == nil && resp.StatusCode != 200 {
		err = fmt.Errorf("HTTP response %s", resp.Status)
	}
	return body, err
}

func (c *Client) login() error {
	body, err := c.get("https://www.opal.com.au/login/index")
	if err != nil {
		return fmt.Errorf("GETting login form: %v", err)
	}
	token, err := parseLogin(body)
	if err != nil {
		return err
	}
	form := url.Values{
		"h_username": []string{c.username},
		"h_password": []string{c.password},
		"CSRFToken":  []string{token},
	}
	resp, err := c.hc.PostForm("https://www.opal.com.au/login/registeredUserUsernameAndPasswordLogin", form)
	if err != nil {
		return fmt.Errorf("POSTing login form: %v", err)
	}
	_, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("reading login form response: %v", err)
	}
	// A successful response sets a cookie in c.hc.
	if resp.StatusCode != 200 {
		return fmt.Errorf("login form response was %s", resp.Status)
	}
	return nil
}
