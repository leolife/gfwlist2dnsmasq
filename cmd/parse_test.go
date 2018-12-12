package cmd

import (
	"strings"
	"testing"
)

func TestGetHostname(t *testing.T) {
	h, err := getHostname("google.com")
	if err != nil {
		t.Fail()
	}
	if h != "google.com" {
		t.Fail()
	}
}

func TestValidIPv4(t *testing.T) {
	ok := isIPv4("127.0.0.1")
	if !ok {
		t.Fail()
	}
}

func TestInvalidIPv4(t *testing.T) {
	ok := isIPv4("google.com")
	if ok {
		t.Fail()
	}
}

func TestReduceDomain(t *testing.T) {
	h := reduceDomain("www.google.com")
	if h != "google.com" {
		t.Fail()
	}
}

func TestParseList(t *testing.T) {
	r := strings.NewReader("www.google.com")
	c, err := ParseList(r)
	if err != nil {
		t.Fail()

	}
	if _, ok := c["google.com"]; !ok {
		t.Fail()
	}
}
