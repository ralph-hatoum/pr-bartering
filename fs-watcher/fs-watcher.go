package fswatcher

import (
	datastructures "bartering/data-structures"
	// "bartering/functions"
	storagerequests "bartering/storage-requests"
	"fmt"
	"log"

	fsnotify "github.com/fsnotify/fsnotify"
)

func FsWatcher(path string, storage_pool []string, peerScores []datastructures.NodeScore, K int, storageRequest datastructures.StorageRequest, port string, bytesAtPeers []datastructures.PeerStorageUse, fulfilledRequests []datastructures.FulfilledRequest, scoreDecreaseRefStoReq float64) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					filePath := event.Name
					// Handle the file creation event
					fmt.Println("New file ", filePath, " detected, storing on network")
					// go functions.Store(filePath, storage_pool, pendingRequests) // this still does not actually trigger storage on the network
					go storagerequests.StoreKCopiesOnNetwork(peerScores, K, storageRequest, port, bytesAtPeers, &fulfilledRequests, scoreDecreaseRefStoReq)
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					filePath := event.Name
					log.Println("Modified file:", filePath)

					// Handle the file modification event
					// basically considered a new file ; so old version should be deleted and new version stored
					// need to implement a delRq message but security for this will be hard to implement
					// maybe we should just rely on not renewing the lease
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					filePath := event.Name
					log.Println("Removed file:", filePath)
					// Handle the file removal event
					// there is no distinction between file removal and file renaming (renaming is basically just deleting and recreating a file)
					// in ipfs it is the same logic, not sure if this is problematic
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error:", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
