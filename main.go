package main

import (
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var fileBatch []string
	if len(os.Args) != 2 {
		os.Exit(3)
	}

	lokasi := os.Args[1]

	err := filepath.Walk(lokasi,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) != "" {
				fileBatch = append(fileBatch, path)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	for index, value := range fileBatch {
		hash, err := hashFileCrc32(value, 0xedb88320)
		if err == nil {
			fmt.Println("index[", index, "] :", value, " : ", hash)
		}
	}
}

func hashFileCrc32(filePath string, polynomial uint32) (string, error) {
	var returnCRC32String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnCRC32String, err
	}
	defer file.Close()
	tablePolynomial := crc32.MakeTable(polynomial)
	hash := crc32.New(tablePolynomial)
	if _, err := io.Copy(hash, file); err != nil {
		return returnCRC32String, err
	}
	hashInBytes := hash.Sum(nil)[:]
	returnCRC32String = hex.EncodeToString(hashInBytes)
	return returnCRC32String, nil

}
