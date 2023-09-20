package main

import (
	"context"
	"fmt"

	"go.formulabun.club/extractor"
	"go.formulabun.club/metadatadb"
	"go.formulabun.club/srb2kart/network"
)

func scrapeServer(host string, ctx context.Context) error {
	serverInfo, _, err := network.AskInfo(host)
	fmt.Println("serverinfo")
	if err != nil {
		return err
	}

	files, err := network.TellAllFilesNeeded(host)
	if err != nil {
		return err
	}
	dbFiles := make([]metadatadb.File, len(files))
	for i, f := range files {
		dbFiles[i] = metadatadb.File{f.Filename, f.Checksum.String()}
	}

	c, err := metadatadb.NewClient(ctx)
	if err != nil {
		return err
	}

	filesToExtract := make(chan metadatadb.File)
	ec := extractor.NewClient()

	steps := make(chan error, 3)
	defer close(steps)

	go func() {
		steps <- errorf("could not store file info: %w", storeFileInfo(c, dbFiles, filesToExtract, ctx))
	}()
	go func() {
		steps <- errorf("could not store server info: %w", storeServerInfo(c, serverInfo, host, dbFiles, ctx))
	}()
	go func() {
		steps <- errorf("could not extract file data: %w", extractFiles(ec, serverInfo, filesToExtract, ctx))
	}()

	return ReadErrors(steps, 3, ctx)
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
