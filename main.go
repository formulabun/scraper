package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.formulabun.club/masterserver"
	"go.formulabun.club/metadatadb"
)

const mainDuration = time.Minute * 10

func mainTicker() <-chan time.Time {
	res := make(chan time.Time)
	go func() {
		res <- time.Now()

		tic := time.Tick(mainDuration)
		for t := range tic {
			res <- t
		}
	}()
	return res
}

func scrapeMs(c *metadatadb.Client) {
	for _ = range mainTicker() {
		servers, _ := masterserver.ListServers()
		serverHosts := make([]string, len(servers))
		for i, s := range servers {
			serverHosts[i] = fmt.Sprintf("%s:%s", s.Ip, s.Port)
		}
		scrapeServers(serverHosts, c)
	}
}

func main() {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Second*5)
	c, err := metadatadb.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	switch len(os.Args) {
	case 1:
		scrapeMs(c)
	case 2:
		scrapeServers([]string{os.Args[1]}, c)
	default:
		fmt.Printf("usage: %s [host?]\n", os.Args[0])
	}

}
