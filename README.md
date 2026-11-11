# dnsexpect

A small Go CLI and library that checks whether a hostname resolves to the IP addresses you expect, and can wait until it does. Standard library only.

## Why

After pointing a domain at a new server, requesting a Let's Encrypt certificate too early fails the HTTP challenge and can hit rate limits. Waiting for DNS to show the new address first avoids that. It is also a quick sanity check in deploy scripts and CI.

## Install

```
go install github.com/PrinceAdvant/dnsexpect/cmd/dnsexpect@latest
```

## Usage

```
dnsexpect -ip 203.0.113.10 app.example.com
dnsexpect -ip 203.0.113.10,2001:db8::10 -wait 15m -interval 30s app.example.com
```

Exit code 0 means every expected IP was returned and nothing else. Exit code 1 means a mismatch or lookup error, with the missing and extra addresses printed. Exit code 2 is a usage error.

As a library:

```go
res, err := dnsexpect.Check(ctx, nil, "app.example.com", []string{"203.0.113.10"})
if err == nil && !res.Matches {
    log.Printf("missing %v extra %v", res.Missing, res.Extra)
}
```

## Notes

The check uses the system resolver, which may cache answers. Other resolvers around the world can still return the old record until its TTL expires. Why that happens is explained in [DNS propagation explained](https://peon.sh/blogs/dns-propagation-explained). Built by a contributor to [Peon](https://peon.sh/), an open-source self-hosted deployment platform.

MIT licensed.
