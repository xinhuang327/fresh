package runner

import (
	"os"
	"path/filepath"
	"strings"

	"encoding/json"
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

const (
	CommitType_TryDeploy     = "try dep"
	CommitType_ReleaseDeploy = "rel dep"
)

func commitTypeForMessage(msg string) string {
	lowerMsg := strings.ToLower(msg)
	if strings.Index(lowerMsg, CommitType_TryDeploy) > -1 {
		return CommitType_TryDeploy
	}
	if strings.Index(lowerMsg, CommitType_ReleaseDeploy) > -1 {
		return CommitType_ReleaseDeploy
	}
	return ""
}

func watchGitHandler(w http.ResponseWriter, r *http.Request) {
	//gitLog("Got request: %s", r.URL.Path) // /drone or /gogs
	isDrone := r.URL.Path == "/drone"
	isGogs := !isDrone
	if r.Method == "POST" {
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			gitLog("Read payload error %s", err.Error())
			w.WriteHeader(500)
			w.Write([]byte("Fresh: Error reading payload: " + err.Error()))
			return
		} else {
			var commitType string = ""
			if isDrone {
				var dronePayload DronePayload
				err = json.Unmarshal(payload, &dronePayload)
				if err != nil {
					w.WriteHeader(500)
					w.Write([]byte("Fresh: Unmarshal DronePayload error: " + err.Error()))
					return
				}
				if dronePayload.Build.Status == "success" {
					gitLog("Commit message: %s", dronePayload.Build.Message)
					commitType = commitTypeForMessage(dronePayload.Build.Message)
				} else {
					gitLog("Drone build failed")
				}
			} else {
				var gogsPayload GogsPayload
				err = json.Unmarshal(payload, &gogsPayload)
				if err != nil {
					w.WriteHeader(500)
					w.Write([]byte("Fresh: Unmarshal GogsPayload error: " + err.Error()))
					return
				}
				if len(gogsPayload.Commits) > 0 {
					gitLog("Commit message: %s", gogsPayload.Commits[0].Message)
					commitType = commitTypeForMessage(gogsPayload.Commits[0].Message)
				}
			}
			payloadStr := string(payload)
			gitLog("Got payload from %s, type: %s", r.URL.Path, commitType)
			if commitType == CommitType_TryDeploy && isGogs ||
				commitType == CommitType_ReleaseDeploy && isDrone {
				gitLog("Start rebuild...")
				gitStartChannel <- payloadStr
			}
			w.Write([]byte("Fresh: OK"))
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
