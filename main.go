package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

var (
	deleteAfter time.Duration
	sleep       time.Duration
	directories []string
)

func clean(dir string) error {
	stats, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	now := time.Now()
	for _, stat := range stats {
		if stat.Mode().IsRegular() && stat.ModTime().Add(deleteAfter).Before(now) {
			fp := path.Join(dir, stat.Name())
			if err := os.Remove(fp); err != nil {
				log.Printf("Failed to delete file %s", fp)
			}
		}
	}
	log.Printf("Successfully deleted %s", dir)
	return nil
}

func env(key string, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func main() {
	var err error
	if deleteAfter, err = time.ParseDuration(env("DELETE_AFTER", "30m")); err != nil {
		log.Fatalf("Invalid DELETE_AFTER provided. See the 'https://golang.org/pkg/time/#ParseDuration'")
	}
	if sleep, err = time.ParseDuration(env("SLEEP", "30m")); err != nil {
		log.Fatalf("Invalid RUN_EVERY provided. See the 'https://golang.org/pkg/time/#ParseDuration'")
	}
	directories = strings.Split(env("DIRS", ""), ",")
	log.Printf("Will run every %v", sleep)
	log.Printf("Delete files that was changed %v ago", deleteAfter)

	for i, dir := range directories {
		dir = strings.TrimSpace(dir)
		stat, err := os.Stat(dir)
		if err != nil {
			log.Fatalf("Failed to stat directory '%s': %v", dir, err)
		}
		if !stat.IsDir() {
			log.Fatalf("Error: '%s' is not a directory", dir)
		}
		directories[i] = dir
	}
	for {
		for _, d := range directories {
			if err := clean(d); err != nil {
				log.Fatalf("Failed to clean directory '%s': %v", d, err)
			}
		}
		time.Sleep(sleep)

	}
}
