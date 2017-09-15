package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type logWriter struct{}

func (writer logWriter) Write(bytes []byte) (int, error) {
	printTime := color.New(color.FgCyan, color.Bold).SprintFunc()
	timestamp := time.Now().Format("2006-01-02 15:04:05 PM")
	return fmt.Printf("%s %s", printTime(timestamp), string(bytes))
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
	excludePatterns := []string{
		".git",
		"node_modules",
		".idea",
		".vagrant",
		".sass-cache",
	}
	err = RecursiveWatch(baseDir, watcher, excludePatterns, true)
	if err != nil {
		log.Fatal(err)
	}

	// Hang so program doesn't exit
	<-done
}

func RecursiveWatch(path string, watcher *fsnotify.Watcher, exclude []string, debug bool) (err error) {
	err = filepath.Walk(path, func(walkPath string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !stringInSlice(walkPath, exclude) {
			if debug {
				log.Printf("Watching: %s", walkPath)
			}
			if err = watcher.Add(walkPath); err != nil {
				return err
			}
		}

		return nil
	})
	return err
}

func stringInSlice(str string, list []string) bool {
	for _, pattern := range list {
		match, err := regexp.MatchString(pattern, str)

		if err != nil {
			log.Fatal("Trouble excluding directory")
		}

		if match {
			return true
		}
	}
	return false
}
