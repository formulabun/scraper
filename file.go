package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"go.formulabun.club/extractor"
	"go.formulabun.club/functional/strings"
	"go.formulabun.club/metadatadb"
	"go.formulabun.club/srb2kart/network"
	"go.formulabun.club/storage"
)

func extractFiles(c *extractor.Client, info network.ServerInfo, files chan metadatadb.File, ctx context.Context) error {
	source, err := url.Parse(strings.SafeNullTerminated(info.HttpSource[:]))
	if err != nil {
		return err
	}

	errs := make([]error, 0)

	// nested function for defer close
	saveFile := func(f metadatadb.File) error {
		resp, err := http.Get(source.JoinPath(f.Filename).String())
		if err != nil || resp.StatusCode/100 != 2 {
			// TODO rm file from db
			return err
		}
		defer resp.Body.Close()

		err = storage.Store(f, resp.Body)
		if err != nil {
			return err
		}
		return nil
	}

	for {
		var err error
		select {
		case <-ctx.Done():
			break
		case f, ok := <-files:
			if !ok {
				break
			}
			fmt.Printf("new file: %s\n", f.Filename)
			if !storage.Has(f) {
				err = saveFile(f)
			}
			if err != nil {
				errs = append(errs, err)
			} else {
				c.ExtractFile(f)
			}
		}
	}

	return errors.Join(errs...)
}

func storeFileInfo(c *metadatadb.Client, files []metadatadb.File, filesToExtract chan metadatadb.File, ctx context.Context) error {
	fileChan := make(chan metadatadb.File)
	go func() {
		defer close(fileChan)
		for _, f := range files {
			select {
			case <-ctx.Done():
				return
			default:
				fileChan <- f
			}
		}
	}()

	defer close(filesToExtract)

	for {
		select {
		case <-ctx.Done():
			return errors.New("Cancelled operation")
		case f, ok := <-fileChan:
			if !ok {
				return nil
			}
			existed, err := c.AddFile(f, ctx)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if !existed {
				filesToExtract <- f
			}
		}
	}
	return nil
}
