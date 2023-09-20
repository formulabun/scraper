package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.formulabun.club/masterserver"
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

func scrapeMs() {
	for _ = range mainTicker() {
		servers, _ := masterserver.ListServers()
		for _, server := range servers {
			loc := fmt.Sprintf("%s:%s", server.Ip, server.Port)
			fmt.Printf("scraping %s\n", loc)
			ctx := context.Background()
			ctx, _ = context.WithTimeout(ctx, time.Second*120)
			err := scrapeServer(loc, ctx)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func scrapeSingle(host string) {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Second*120)
	err := scrapeServer(host, ctx)
	if err != nil {
		panic(err)
	}
}

func main() {
	switch len(os.Args) {
	case 1:
		scrapeMs()
	case 2:
		scrapeSingle(os.Args[1])
	default:
		fmt.Printf("usage: %s [host?]\n", os.Args[0])
	}

}
