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
		for _, server := range servers {
			loc := fmt.Sprintf("%s:%s", server.Ip, server.Port)
			fmt.Printf("scraping %s\n", loc)
			ctx := context.Background()
			ctx, _ = context.WithTimeout(ctx, time.Minute)
			err := scrapeServer(loc, c, ctx)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func scrapeSingle(host string, c *metadatadb.Client) {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Millisecond*200)
	err := scrapeServer(host, c, ctx)
	if err != nil {
		panic(err)
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
		scrapeSingle(os.Args[1], c)
	default:
		fmt.Printf("usage: %s [host?]\n", os.Args[0])
	}

}
