package main

import (
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var fileBatch []string
	if len(os.Args) < 2 {
		os.Exit(3)
	}

	//TODO tambah supaya bisa multi input
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
		fmt.Println("Please use double quetes (\") while input path")
		fmt.Println("ex: \"d:\\test\\test.go\"")
	}

	for index, value := range fileBatch {
		hash, err := hashFileCrc32(value, 0xedb88320)
		if err == nil {
			fmt.Println("index[", index, "] :", value, " : ", strings.ToUpper(hash))
			path := filepath.Dir(value)
			file := filepath.Base(value)
			format := filepath.Ext(value)
			filename := file[0 : len(file)-len(format)]
			fmt.Println(value, path+"\\"+filename+" ["+hash+"]"+format)
			//os.Rename(value, path+"\\"+filename+" ["+hash+"]"+format)
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
