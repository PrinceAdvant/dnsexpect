package dnsexpect

import (
	"context"
	"net"
	"testing"
	"time"
)

type fake struct{ ips []string }

func (f fake) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	var out []net.IPAddr
	for _, s := range f.ips {
		out = append(out, net.IPAddr{IP: net.ParseIP(s)})
	}
	return out, nil
}

func TestMatch(t *testing.T) {
	r, err := Check(context.Background(), fake{[]string{"203.0.113.10"}}, "app.example.com", []string{"203.0.113.10"})
	if err != nil || !r.Matches {
		t.Fatalf("%+v %v", r, err)
	}
}

func TestMismatch(t *testing.T) {
	r, _ := Check(context.Background(), fake{[]string{"198.51.100.7"}}, "app.example.com", []string{"203.0.113.10"})
	if r.Matches || len(r.Missing) != 1 || len(r.Extra) != 1 {
		t.Fatalf("%+v", r)
	}
}

func TestInvalidWant(t *testing.T) {
	if _, err := Check(context.Background(), fake{}, "x", []string{"nope"}); err == nil {
		t.Fatal("expected error")
	}
}

func TestWaitTimeout(t *testing.T) {
	_, err := Wait(context.Background(), fake{[]string{"198.51.100.7"}}, "x", []string{"203.0.113.10"}, 10*time.Millisecond, 30*time.Millisecond)
	if err == nil {
		t.Fatal("expected timeout")
	}
}
