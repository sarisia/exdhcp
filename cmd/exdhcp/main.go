package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"github.com/sarisia/exdhcp"
)

func main() {
	parent, cancel := context.WithCancel(context.Background())
	defer cancel()

	// arg parsing
	ifname := flag.String("if", "", "Interface name")
	serverAddr := flag.String("server", "", "DHCP server IP")
	verbose := flag.Bool("verbose", false, "Print debug logs")
	numTries := flag.Int("num", 10, "Number of tries. Set 0 to do exhaustion attack.")
	timeout := flag.Int("timeout", 10, "Timeout seconds to wait DHCP offer")
	release := flag.Bool("release", true, "Release DHCP lease")
	exportCSV := flag.Bool("csv", true, "Export leases to CSV file")

	flag.Parse()

	// signal
	ctx, _ := signal.NotifyContext(parent, os.Interrupt)

	cli, err := exdhcp.New(*ifname, *serverAddr, *verbose)
	if err != nil {
		return
	}

	cli.Start(ctx, *numTries, *timeout, *release, *exportCSV)
}
