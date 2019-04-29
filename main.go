package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	var fileBatch []string
	var fileNewName []string
	if len(os.Args) < 2 {
		os.Args = append(os.Args, interfaces())
	}

	for index, value := range os.Args {

		if index == 0 {
			continue
		}

		lokasi := value

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
	}

	askForRename := false
	for index, value := range fileBatch {
		hash, err := hashFileCrc32(value, 0xedb88320)
		hash = strings.ToUpper(hash)
		if err == nil {
			path := filepath.Dir(value)
			file := filepath.Base(value)
			format := filepath.Ext(value)
			filename := file[0 : len(file)-len(format)]

			if findCrc(value, hash) == "File OK" {
				newName := ""
				fmt.Printf("[%d] %s %s\n", index, "[   File OK   ]", value)
				fileNewName = append(fileNewName, newName)
			} else if findCrc(value, hash) == "CRC Not Found" {
				newName := path + "\\" + filename + " [" + hash + "]" + format
				fmt.Printf("[%d] %s %s -> %s\n", index, "[CRC Not Found]", value, newName)
				fileNewName = append(fileNewName, newName)
				askForRename = true
			} else {
				newName := ""
				fmt.Printf("[%d] %s %s %s\n", index, "[File Corrupt ]", value, hash)
				fileNewName = append(fileNewName, newName)
			}
		}
	}
	//TODO buat agar bisa select rename file
	//TODO buat supaya bisa di sort by category
	if askForRename {
		fmt.Print("Rename all CRC Not Found(y/n) ?")
		if askForConfirmation() {
			renamer(fileBatch, fileNewName)
			pressAnyKey()
		}
	} else {
		pressAnyKey()
	}
}

func hashFileCrc32(filePath string, polynomial uint32) (string, error) {
	var returnCRC32String string
	fmt.Println("Now Checking :", filePath, "...")
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
	var validID = regexp.MustCompile(`\[([[:xdigit:]]{8})\]`)
	crc := string(validID.Find([]byte(filename)))
	crc = strings.ToUpper(crc)
	if crc == "["+hash+"]" {
		return "File OK"
	} else if crc == "" {
		return "CRC Not Found"
	} else {
		return "File Corrupt"
	}
}

func renamer(origin []string, newName []string) {
	if len(origin) != len(newName) {
		fmt.Println("Panjang tidak sama, BUG!!!!")
		return
	} else {
		for index, value := range newName {
			if value != "" {
				os.Rename(origin[index], value)
			}
		}
		fmt.Println("All Done")
	}
}

func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return askForConfirmation()
	}
}

func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}

func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

func pressAnyKey() {
	var pressAnyKey string
	fmt.Println("Press Any Key to Continue")
	n, _ := fmt.Scanln(&pressAnyKey)
	_ = n
}

func interfaces() string {

	reader := bufio.NewReader(os.Stdin)

	fmt.Println(
		`===============================
           RunMe v1.1
================================`)
	fmt.Print("Enter paths : ")
	patssshs, _ := reader.ReadString('\n')
	//paths = strings.TrimSuffix(paths, "\n")
	patssshs = patssshs[0:(len(patssshs) - 2)]

	return patssshs
}
