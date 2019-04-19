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
	var fileNewName []string
	if len(os.Args) < 2 {
		fmt.Println("Need location of file")
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

	for _, value := range fileBatch {
		hash, err := hashFileCrc32(value, 0xedb88320)
		if err == nil {
			//fmt.Print("index[", index, "] : ", value)
			path := filepath.Dir(value)
			file := filepath.Base(value)
			format := filepath.Ext(value)
			filename := file[0 : len(file)-len(format)]
			newName := path + "\\" + filename + " [" + hash + "]" + format
			findCrc(value, hash)
			//fmt.Println(value, path+"\\"+filename+" ["+hash+"]"+format)
			//fmt.Println(" ->", newName)
			fileNewName = append(fileNewName, newName)
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

func findCrc(filename string, hash string) string {
	//fmt.Print(strings.Trim("[Hello], Gophers!!!", "[,]"))
	fmt.Print(strings.TrimRight(strings.TrimLeft(" asd [Hello], Gophers!!!", "]["), "]["))
	return "File OK"
	//return "No CRC Found"
	//return "File Corrupt"
}
