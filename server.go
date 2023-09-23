package main

import (
	"context"
	"fmt"

	"go.formulabun.club/extractor"
	"go.formulabun.club/metadatadb"
	"go.formulabun.club/srb2kart/network"
	bunstrings "go.formulabun.club/srb2kart/strings"
)

func scrapeServer(host string, c *metadatadb.Client, ctx context.Context) error {
	serverInfo, _, err := network.AskInfo(host)
	if err != nil {
		return err
	}
	fmt.Printf("starting %s\n", bunstrings.RemoveColorCodes(string(serverInfo.ServerName[:])))

	files, err := network.TellAllFilesNeeded(host)
	if err != nil {
		return err
	}
	fmt.Printf("%v files\n", len(files))

	dbFiles := make([]metadatadb.File, len(files))
	for i, f := range files {
		dbFiles[i] = metadatadb.File{f.Filename, f.Checksum.String()}
	}

	filesToExtract := make(chan metadatadb.File, 5)
	ec := extractor.NewClient()

	doneChan := make(chan struct{})

	steps := make(chan error, 3)
	defer close(steps)

	go func() {
		steps <- errorf("could not store file info: %w", storeFileInfo(c, dbFiles, filesToExtract, doneChan, ctx))
	}()
	go func() {
		steps <- errorf("could not store server info: %w", storeServerInfo(c, serverInfo, host, dbFiles, ctx))
	}()
	go func() {
		steps <- errorf("could not extract file data: %w", extractFiles(ec, serverInfo, filesToExtract, doneChan, ctx))
	}()

	err = ReadErrors(steps, 3, ctx)
	fmt.Println()
	return err
}

func storeServerInfo(c *metadatadb.Client, info network.ServerInfo, host string, files []metadatadb.File, ctx context.Context) error {
	steps := make(chan error, 2)
	defer close(steps)

	go func() {
		steps <- c.AddServerInfo(host, info, ctx)
	}()

	go func() {
		steps <- c.AddServerFiles(host, files, ctx)
	}()

	return ReadErrors(steps, 2, ctx)
}
