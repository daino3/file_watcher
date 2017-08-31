package file_watcher

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"log"
	"os"
	"regexp"
	"bufio"
)

type VagrantFile struct {
	fpath string
	valid bool
}

func (v *VagrantFile) parse() []string {
	file, err := os.Open(v.fpath)
	defer file.Close()

	if err != nil {
		fmt.Printf("error opening Vagrant file: %v\n", err)
		os.Exit(1)
	}

	re_is_folder := regexp.MustCompile(`\s?config.vm.synced_folder\s?`)
	re_is_vagrant := regexp.MustCompile(`\/vagrant`)

	var lines []string
	scanner := bufio.NewScanner(file)
    for scanner.Scan() {
		if text := scanner.Text(); re_is_folder.MatchString(text) && re_is_vagrant.MatchString(text) {
			lines = append(lines, text)
			fmt.Println("Matched line: %s", text)
		} else {
			fmt.Println("Did Not Matched line: %s", text)
		}
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

	return lines
}

func main() {
	fmt.Printf("hello, world\n")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				log.Println("event:", ev)
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	vfile := VagrantFile{fpath: "./VagrantFile"}

	for _, file := range vfile.parse() {
		os.Expand()

	    err = watcher.AddWatch(file)
	    if err != nil {
	        log.Fatal(err)
	    }
	}

	err = watcher.Watch("testDir")
	if err != nil {
		log.Fatal(err)
	}

	// Hang so program doesn't exit
	<-done

	/* ... do stuff ... */
	watcher.Close()
}
