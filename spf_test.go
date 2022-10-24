package spf

import (
	"context"
	"testing"
	"time"
)

const domain = "google.com"

type spferror struct {
	domain string
	raw    string
}

type spftest struct {
	server string
	email  string
	result Result
}

type spfstr struct {
	raw      string
	expected string
}

func TestNewSPF(t *testing.T) {
	errorTests := []spferror{
		{"google.com", "somestring"},
		{"google.com", "v=spf1 include:_spf.google.com ~all -none"},
		{"google.com", "v=spf1 include:google.com"},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	for _, expected := range errorTests {
		_, err := NewSPF(ctx, expected.domain, expected.raw, 0)

		if err == nil {
			t.Log("Analyzing:", expected.raw)
			t.Error("Expected error got nil")
		}
	}
}

func TestSPFTest(t *testing.T) {
	tests := []spftest{
		{"127.0.0.1", "info@google.com", SoftFail},
		{"74.125.141.26", "info@google.com", Pass},
		{"35.190.247.0", "info@google.com", Pass},
		{"172.217.0.0", "info@_netblocks3.google.com", Pass},
		{"172.217.0.0", "info@google.com", Pass},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	for _, expected := range tests {
		actual, err := SPFTest(ctx, expected.server, expected.email)
		if err != nil {
			t.Error(err)
		}

		if actual != expected.result {
			t.Error("For", expected.server, "at", expected.email, "Expected", expected.result, "got", actual)
		}
	}
}

func TestSPFString(t *testing.T) {
	tests := []spfstr{
		{
			"v=spf1 ip4:45.55.100.54 ip4:192.241.161.190 ip4:188.226.145.26 ~all",
			"v=spf1 ip4:45.55.100.54 ip4:192.241.161.190 ip4:188.226.145.26 ~all",
		},
		{
			"v=spf1 ip4:127.0.0.0/8 -ip4:127.0.0.1 ?ip4:127.0.0.2 -all",
			"v=spf1 ip4:127.0.0.0/8 -ip4:127.0.0.1 ?ip4:127.0.0.2 -all",
		},
		{
			"v=spf1 redirect=_spf.sample.invalid",
			"v=spf1 redirect=_spf.sample.invalid",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	for _, tcase := range tests {
		s, err := NewSPF(ctx, "domain", tcase.raw, 0)
		if err != nil {
			t.Log("Analyzing", tcase.raw)
			t.Error(err)
		}

		r := s.SPFString()
		if r != tcase.expected {
			t.Log("Analyzing", tcase.raw)
			t.Error("Expected", tcase.expected, "got", r)
		}
	}
}
