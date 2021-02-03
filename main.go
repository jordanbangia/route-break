package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aeden/traceroute"
	"github.com/go-ping/ping"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	go func() {
		<-c
		cancel()
	}()

	fmt.Println("starting up")
	for {
		addresses, err := getAddressesInRouteTo("8.8.8.8")
		if err != nil {
			fmt.Printf("error running traceroute: %s\n", err)
			return
		}

		group, groupCtx := errgroup.WithContext(ctx)
		for _, address := range addresses {
			group.Go(func() error {
				return pingAddress(groupCtx, address)
			})
		}

		err = group.Wait()
		if err != nil && ctx.Err() == nil {
			fmt.Println("===== route is broken =====")
			for i, address := range addresses {
				fmt.Printf("%d: %s\n", i, address)
			}

			fmt.Printf("route is broken: %s\n", err)
		} else if ctx.Err() != nil {
			return
		}
	}
}

func getAddressesInRouteTo(routeRoot string) ([]string, error) {
	routeResult, err := traceroute.Traceroute(routeRoot, &traceroute.TracerouteOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to run route to %s: %w", routeRoot, err)
	}

	addresses := []string{}
	for _, hop := range routeResult.Hops {
		addresses = append(addresses, hop.AddressString())
	}
	return addresses, nil
}

func pingAddress(ctx context.Context, address string) error {
	pinger, err := ping.NewPinger(address)
	if err != nil {
		return fmt.Errorf("failed to create pinger to %s: %w", address, err)
	}

	pinger.SetPrivileged(true)

	// do 10 pings, then call it a wrap
	pinger.Count = 10
	err = pinger.Run()

	if err != nil {
		return fmt.Errorf("error pinging %s: %w", address, err)
	}

	stats := pinger.Statistics()

	if stats.PacketLoss > 0.55 {
		return fmt.Errorf("lost too many packets pinging %s", address)
	}
	return ctx.Err()
}
