// Command dnsexpect exits 0 once a hostname resolves to the expected IPs.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PrinceAdvant/dnsexpect"
)

func main() {
	want := flag.String("ip", "", "comma separated list of expected IP addresses")
	wait := flag.Duration("wait", 0, "keep checking for this long, for example 10m")
	interval := flag.Duration("interval", 15*time.Second, "delay between checks when -wait is set")
	flag.Parse()
	if flag.NArg() != 1 || *want == "" {
		fmt.Fprintln(os.Stderr, "usage: dnsexpect -ip 203.0.113.10[,2001:db8::1] [-wait 10m] host")
		os.Exit(2)
	}
	ips := strings.Split(*want, ",")
	ctx := context.Background()
	var res dnsexpect.Result
	var err error
	if *wait > 0 {
		res, err = dnsexpect.Wait(ctx, nil, flag.Arg(0), ips, *interval, *wait)
	} else {
		res, err = dnsexpect.Check(ctx, nil, flag.Arg(0), ips)
		if err == nil && !res.Matches {
			err = fmt.Errorf("missing %v, extra %v", res.Missing, res.Extra)
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("%s resolves to %v\n", res.Host, res.Got)
}
