package fswatcher

import (
	"log"

	fsnotify "github.com/fsnotify/fsnotify"
)

func FsWatcher(path string) {
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
					log.Println("Created file:", event.Name)
					// Handle the file creation event
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("Modified file:", event.Name)
					// Handle the file modification event
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					log.Println("Removed file:", event.Name)
					// Handle the file removal event
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
