package runner

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	startChannel    chan string
	gitStartChannel chan string
	stopChannel     chan bool
	mainLog         logFunc
	watcherLog      logFunc
	runnerLog       logFunc
	buildLog        logFunc
	appLog          logFunc
	gitLog          logFunc
)

func flushEvents() {
	for {
		select {
		case eventName := <-startChannel:
			mainLog("receiving event %s", eventName)
		default:
			return
		}
	}
}

func flushGitEvents() {
	for {
		select {
		case gitInfo := <-gitStartChannel:
			mainLog("receiving gitInfo %s", gitInfo)
		default:
			return
		}
	}
}

func gitStart() {
	loopIndex := 0
	buildDelay := buildDelay()

	started := false

	go func() {
		for {
			loopIndex++
			mainLog("Waiting (loop %d)...", loopIndex)
			commitType := <-gitStartChannel

			//mainLog("receiving git event")
			//mainLog("receiving first event %s", gitInfo)
			mainLog("sleeping for %d milliseconds", buildDelay)
			time.Sleep(buildDelay * time.Millisecond)
			mainLog("flushing git events")

			flushGitEvents()

			mainLog("Started! (%d Goroutines)", runtime.NumGoroutine())
			err := removeBuildErrorsLog()
			if err != nil {
				mainLog(err.Error())
			}

			// execute git pull
			gitLog("Execute git pull")
			gitCmd := exec.Command("git", "pull")
			output, err := gitCmd.CombinedOutput()
			gitLog("Git pull output:")
			gitLog(string(output))
			if err != nil {
				mainLog(err.Error())
			}

			if commitType == CommitType_UpdateResources && started {
				// restart without build
				stopChannel <- true
				run()
			} else if !alreadyHasBuild() {
				errorMessage, ok := build()
				if !ok {
					mainLog("Build Failed: \n %s", errorMessage)
					if !started {
						//os.Exit(1)
						continue
					}
					createBuildErrorsLog(errorMessage)
				} else {
					writeBuildRevision()
					if started {
						stopChannel <- true
					}
					run()
				}
			} else {
				if started {
					stopChannel <- true
				}
				run()
			}

			started = true
			mainLog(strings.Repeat("-", 20))
		}
	}()
}

func start() {
	loopIndex := 0
	buildDelay := buildDelay()

	started := false

	go func() {
		for {
			loopIndex++
			mainLog("Waiting (loop %d)...", loopIndex)
			eventName := <-startChannel

			mainLog("receiving first event %s", eventName)
			mainLog("sleeping for %d milliseconds", buildDelay)
			time.Sleep(buildDelay * time.Millisecond)
			mainLog("flushing events")

			flushEvents()

			mainLog("Started! (%d Goroutines)", runtime.NumGoroutine())
			err := removeBuildErrorsLog()
			if err != nil {
				mainLog(err.Error())
			}

			// use "go run" instead of build and run
			if runFile() != "" {
				if started {
					stopChannel <- true
				}
				run()
			} else {
				errorMessage, ok := build()
				if !ok {
					mainLog("Build Failed: \n %s", errorMessage)
					if !started {
						os.Exit(1)
					}
					createBuildErrorsLog(errorMessage)
				} else {
					if started {
						stopChannel <- true
					}
					run()
				}
			}
			started = true
			mainLog(strings.Repeat("-", 20))
		}
	}()
}

const VersionFileName = "build_revision.txt"

func alreadyHasBuild() bool {
	_, err := os.Stat(buildPath())
	noExistingBuild := err != nil
	if noExistingBuild {
		return false
	}
	mainLog("Check if already has a build...")
	// git rev-parse HEAD
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		mainLog("Get git revision error: %s", err)
		return false
	}
	filePath := filepath.Join(tmpPath(), VersionFileName)
	existingRev, err := ioutil.ReadFile(filePath)
	if err != nil {
		mainLog("Get build revision error: %s", err)
		return false
	}
	gitRev := string(output)
	buildRev := string(existingRev)
	mainLog("git rev: %s, build rev: %s", gitRev, buildRev)
	if gitRev == buildRev {
		mainLog("Revision equals. Skip build.")
	}
	return gitRev == buildRev
}

func writeBuildRevision() {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		mainLog("Get git revision error", err)
		return
	}
	filePath := filepath.Join(tmpPath(), VersionFileName)
	err = ioutil.WriteFile(filePath, output, 0664)
	if err != nil {
		mainLog("Write revision file error", err)
		return
	}
	mainLog("writeBuildRevision", string(output))
}

func init() {
	startChannel = make(chan string, 1000)
	gitStartChannel = make(chan string, 1000)
	stopChannel = make(chan bool)
}

func initLogFuncs() {
	mainLog = newLogFunc("main")
	watcherLog = newLogFunc("watcher")
	runnerLog = newLogFunc("runner")
	buildLog = newLogFunc("build")
	appLog = newLogFunc("app")
	gitLog = newLogFunc("git")
}

func setEnvVars() {
	os.Setenv("DEV_RUNNER", "1")
	wd, err := os.Getwd()
	if err == nil {
		os.Setenv("RUNNER_WD", wd)
	}

	for k, v := range settings {
		key := strings.ToUpper(fmt.Sprintf("%s%s", envSettingsPrefix, k))
		os.Setenv(key, v)
	}
}

// Watches for file changes in the root directory.
// After each file system event it builds and (re)starts the application.
func Start() {
	initLimit()
	initSettings()
	initLogFuncs()
	initFolders()
	setEnvVars()

	_, err := os.Stat(buildPath())
	noExistingBuild := err != nil
	if noExistingBuild {
		mainLog("No existing build, will start build")
	}

	if git_pull_mode() {
		watchGit()
		gitStart()
		gitStartChannel <- "/" // gitStart() will check if already has a build
	} else {
		watch()
		start()
		if noExistingBuild {
			startChannel <- "/"
		}
	}

	<-make(chan int)
}
