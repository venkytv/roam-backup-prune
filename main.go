package main

import (
	"bytes"
	"crypto/sha256"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

const (
	BackupPath   = "Library/Application Support/Roam Research/backups/Duh-Uh"
	BackupPrefix = "backup-Duh-Uh-"
)

func checksum(file string) []byte {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		log.Fatalf("Error calculating checksum: %s: %v", file, err)
	}

	return hasher.Sum(nil)
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	backupDir := path.Join(homeDir, BackupPath)

	backups, err := ioutil.ReadDir(backupDir)
	if err != nil {
		log.Fatal(err)
	}

	lastFileSize := int64(-1)
	var lastChecksum []byte

	for _, backup := range backups {
		backupName := backup.Name()
		if !strings.HasPrefix(backupName, BackupPrefix) {
			log.Print("Ignoring ", backupName)
			continue
		}
		fileSize := backup.Size()
		backupFile := path.Join(backupDir, backupName)
		sum := checksum(backupFile)

		if fileSize == lastFileSize {
			if bytes.Compare(sum, lastChecksum) == 0 {
				log.Printf("Removing: %s: %x", backupName, sum)
				if err := os.Remove(backupFile); err != nil {
					log.Printf("Error removing file: %s: %v", backupFile, err)
				}
				continue
			}
		}
		lastFileSize = backup.Size()
		lastChecksum = sum
		log.Printf("Keeping : %s: %x <-", backupName, lastChecksum)
	}
}
