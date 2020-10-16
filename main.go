package main

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"
)

var (
	deleteAfter time.Duration
	sleep       time.Duration
	directories []string
)

func clean(dir string) error {
	dirLog := logrus.WithField("dir", dir)

	stats, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	now := time.Now()
	cntDeleted := 0
	for _, stat := range stats {
		if stat.Mode().IsRegular() && stat.ModTime().Add(deleteAfter).Before(now) {
			fp := path.Join(dir, stat.Name())
			if err := os.Remove(fp); err != nil {
				dirLog.Errorf("Failed to delete file %s", fp)
			} else {
				cntDeleted += 1
			}
		}
	}
	dirLog.Infof("Cleanup successful, %d files deleted", cntDeleted)
	return nil
}

func env(key string, defaultValue string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return defaultValue
}

func main() {
	var err error
	if deleteAfter, err = time.ParseDuration(env("DELETE_AFTER", "30m")); err != nil {
		logrus.Fatalf("Invalid DELETE_AFTER provided. See the 'https://golang.org/pkg/time/#ParseDuration'")
	}
	if sleep, err = time.ParseDuration(env("SLEEP", "30m")); err != nil {
		logrus.Fatalf("Invalid RUN_EVERY provided. See the 'https://golang.org/pkg/time/#ParseDuration'")
	}
	directories = strings.Split(env("DIRS", ""), ",")
	logrus.Infof("Running the cleanup every %v", sleep)
	logrus.Infof("Deleting files that were changed %v ago", deleteAfter)

	for i, dir := range directories {
		dir = strings.TrimSpace(dir)
		stat, err := os.Stat(dir)
		if err != nil {
			logrus.Fatalf("Failed to stat directory '%s': %v", dir, err)
		}
		if !stat.IsDir() {
			logrus.Fatalf("Error: '%s' is not a directory", dir)
		}
		directories[i] = dir
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

loop:
	for {
		for _, d := range directories {
			if err := clean(d); err != nil {
				log.Fatalf("Failed to clean directory '%s': %v", d, err)
			}
		}
		ticker := time.NewTicker(sleep)
		select {
		case <-c:
			logrus.Infof("Shutting down")
			break loop
		case <-ticker.C:
			ticker.Stop()
		}
	}
}
