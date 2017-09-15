package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	printTime := color.New(color.FgCyan, color.Bold).SprintFunc()
	timestamp := time.Now().Format("2006-01-02 15:04:05 PM")
	return fmt.Printf("%s TRACKING... %s", printTime(timestamp), string(bytes))
}

func main() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	baseDir, err := filepath.Abs(filepath.Dir("./"))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Tracking changes to %s\n", baseDir)

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	done := make(chan bool)

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				log.Println("[EVENT]: ", ev)
			case err := <-watcher.Errors:
				log.Println("[ERROR]: ", err)
			}
		}
	}()

	contents, err := ioutil.ReadDir(baseDir)

	if err != nil {
		log.Fatal(err)
	}

	for _, fileOrDir := range contents {
		err = watcher.Add(baseDir + "/" + fileOrDir.Name())
		if err != nil {
			log.Fatal(err)
		}
	}

	// Hang so program doesn't exit
	<-done

	/* ... do stuff ... */
}
