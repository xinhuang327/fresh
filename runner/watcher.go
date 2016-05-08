package runner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/howeyc/fsnotify"
	"io/ioutil"
	"net/http"
)

func watchFolder(path string) {
	if isIgnoredFolder(path) {
		watcherLog("Ignoring %s", path)
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fatal(err)
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if isWatchedFile(ev.Name) {
					watcherLog("sending event %s", ev)
					startChannel <- ev.String()
				}
			case err := <-watcher.Error:
				watcherLog("error: %s", err)
			}
		}
	}()

	watcherLog("Watching %s", path)
	err = watcher.Watch(path)

	if err != nil {
		fatal(err)
	}
}

func watchGitHandler(w http.ResponseWriter, r *http.Request) {
	gitLog("Got request: ", r.URL.Path)
	if r.Method == "POST" {
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			gitLog("Read payload error %s", err.Error())
			w.WriteHeader(500)
			w.Write([]byte("Fresh: Error reading payload: " + err.Error()))
		} else {
			payloadStr := string(payload)
			gitLog("Got payload:")
			gitLog(payloadStr)
			w.Write([]byte("Fresh: OK"))
			gitStartChannel <- payloadStr
		}
	} else {
		w.Write([]byte("Fresh: running for " + app_name()))
	}
}

func watchGitServer() {
	http.HandleFunc("/drone", watchGitHandler)
	http.HandleFunc("/gogs", watchGitHandler)
	err := http.ListenAndServe(":"+server_port(), nil)
	if err != nil {
		panic("Fresh server error, " + err.Error())
	}
}

func watchGit() {
	go watchGitServer()
}

func watch() {
	root := root()
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && !isTmpDir(path) {
			if len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".") {
				return filepath.SkipDir
			}

			watchFolder(path)
		}

		return err
	})
}
