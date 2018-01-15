package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		log.Println("create dir", dir)
	}

	return nil
}
