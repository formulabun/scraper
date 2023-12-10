package main

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"go.formulabun.club/functional/strings"
	"go.formulabun.club/metadatadb"
	"go.formulabun.club/srb2kart/network"
	"go.formulabun.club/storage"
)

type Downloader struct {
	seenFile map[network.File]struct{}
	newFile  map[network.File]struct{}

	p         chan struct{}
	WaitGroup sync.WaitGroup
}

func newDownloader() Downloader {
	return Downloader{
		make(map[network.File]struct{}),
		make(map[network.File]struct{}),
		make(chan struct{}, 5),
		sync.WaitGroup{},
	}
}

func (d *Downloader) report() {
	fmt.Printf("Seen %d unique files, of which %d were new\n", len(d.seenFile), len(d.newFile))
}

func (d *Downloader) handleFile(f *fileFromServer) (dbFile metadatadb.File) {
	dbFile = metadatadb.File{
		f.file.Filename,
		f.file.Checksum.String(),
	}
	_, seen := d.seenFile[f.file]
	if seen {
		return
	}

	d.WaitGroup.Add(1)
	d.seenFile[f.file] = struct{}{}
	if !storage.Has(dbFile) {
		go func() {
			d.downloadFile(&dbFile, f)
			d.WaitGroup.Done()
		}()
	} else {
		d.WaitGroup.Done()
	}

	return
}

func (d *Downloader) downloadFile(dbFile *metadatadb.File, f *fileFromServer) error {
	d.p <- struct{}{}
	defer func() { <-d.p }()

	fileLocation, err := url.JoinPath(strings.SafeNullTerminated(f.server.HttpSource[:]), f.file.Filename)
	if err != nil {
		return err
	}

	response, err := http.Get(fileLocation)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	storage.Store(*dbFile, response.Body)
	d.newFile[f.file] = struct{}{}
	return nil
}
