// Package dnsexpect checks that a hostname resolves to the IP addresses you expect.
package dnsexpect

import (
	"context"
	"fmt"
	"net"
	"sort"
	"time"
)

// Resolver is the subset of net.Resolver used here, so tests can fake it.
type Resolver interface {
	LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error)
}

// Result compares what a name resolved to with what was expected.
type Result struct {
	Host    string
	Got     []string
	Missing []string
	Extra   []string
	Matches bool
}

// Check resolves host and compares the answer with want.
// The result matches when every wanted IP is present and nothing else is.
func Check(ctx context.Context, r Resolver, host string, want []string) (Result, error) {
	if r == nil {
		r = net.DefaultResolver
	}
	addrs, err := r.LookupIPAddr(ctx, host)
	if err != nil {
		return Result{Host: host}, fmt.Errorf("lookup %s: %w", host, err)
	}
	got := map[string]bool{}
	for _, a := range addrs {
		got[a.IP.String()] = true
	}
	wantSet := map[string]bool{}
	for _, w := range want {
		ip := net.ParseIP(w)
		if ip == nil {
			return Result{Host: host}, fmt.Errorf("invalid expected IP %q", w)
		}
		wantSet[ip.String()] = true
	}
	res := Result{Host: host}
	for ip := range got {
		res.Got = append(res.Got, ip)
		if !wantSet[ip] {
			res.Extra = append(res.Extra, ip)
		}
	}
	for ip := range wantSet {
		if !got[ip] {
			res.Missing = append(res.Missing, ip)
		}
	}
	sort.Strings(res.Got)
	sort.Strings(res.Missing)
	sort.Strings(res.Extra)
	res.Matches = len(res.Missing) == 0 && len(res.Extra) == 0
	return res, nil
}

// Wait polls Check until the answer matches or the timeout elapses.
func Wait(ctx context.Context, r Resolver, host string, want []string, interval, timeout time.Duration) (Result, error) {
	deadline := time.Now().Add(timeout)
	for {
		res, err := Check(ctx, r, host, want)
		if err == nil && res.Matches {
			return res, nil
		}
		if time.Now().After(deadline) {
			if err != nil {
				return res, err
			}
			return res, fmt.Errorf("%s still resolves to %v, want %v", host, res.Got, want)
		}
		select {
		case <-ctx.Done():
			return res, ctx.Err()
		case <-time.After(interval):
		}
	}
}
