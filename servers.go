package main

import (
	"context"
	"fmt"
	"time"

	"go.formulabun.club/extractor"
	"go.formulabun.club/metadatadb"
)

func scrapeServers(hosts []string, c *metadatadb.Client) {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Millisecond*200)

	// TODO sort hosts on distance

	parallel := 5
	filess := make([][]fileFromServer, 0, len(hosts))
	serversChan := make(chan string, parallel)
	filesChan := make(chan []fileFromServer, parallel)

	// This will consume all the servers and put the result in files asap
	go serversToFiles(serversChan, filesChan)

	go func() {
		for _, h := range hosts {
			serversChan <- h
		}
	}()

	for len(filess) != len(hosts) {
		filess = append(filess, <-filesChan)
	}

	close(serversChan)
	close(filesChan)

	oneByOne := make(chan *fileFromServer, 5)
	done := make(chan struct{})
	go transpose(filess, oneByOne, done)

	downloadFiles(oneByOne, done, c)
}

func downloadFiles(oneByOne chan *fileFromServer, done chan struct{}, c *metadatadb.Client) {
	dw := newDownloader()
	ctx := context.Background()
	ec := extractor.NewClient()

	for {
		select {
		case <-done:
			dw.WaitGroup.Wait()
			dw.report()
			return
		case f := <-oneByOne:
			dbFile := dw.handleFile(f)
			err := ec.ExtractFile(dbFile)
			if err != nil {
				fmt.Println(err)
			}
			c.AddFile(dbFile, ctx)
		}
	}
}
