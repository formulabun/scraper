package main

import (
	"errors"
	"log"

	"go.formulabun.club/functional/array"
	"go.formulabun.club/functional/strings"
	"go.formulabun.club/srb2kart/network"
)

type fileFromServer struct {
	file   network.File
	server *network.ServerInfo
}

func serversToFiles(serversIn chan string, filesOut chan []fileFromServer) {
	for server := range serversIn {
		localServer := server
		go func() {
			serverInfo, files, err := getServerData(localServer)
			if err != nil {
				log.Printf("Could not get files for %s: %s\n", localServer, err)
				filesOut <- []fileFromServer{}
			} else {
				serverName := strings.SafeNullTerminated(serverInfo.ServerName[:])
				log.Printf("Got the files for %s\n", serverName)
				filesOut <- array.Map(files, func(f network.File) fileFromServer {
					return fileFromServer{f, &serverInfo}
				})
			}
		}()
	}
}

func getServerData(server string) (network.ServerInfo, []network.File, error) {
	// serverinfo
	type serverInfoErr struct {
		info network.ServerInfo
		err  error
	}
	serverInfoChan := make(chan serverInfoErr)
	go func() {
		serverInfo, _, err := network.AskInfo(server)
		serverInfoChan <- serverInfoErr{serverInfo, err}
	}()
	// files
	type filesErr struct {
		files []network.File
		err   error
	}
	filesChan := make(chan filesErr)
	go func() {
		files, err := network.TellAllFilesNeeded(server)
		filesChan <- filesErr{files, err}
	}()

	infoErr := <-serverInfoChan
	fileErr := <-filesChan

	return infoErr.info, fileErr.files, errors.Join(infoErr.err, fileErr.err)
}
