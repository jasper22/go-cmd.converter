package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
)

func main() {
	pathToExe, _ := os.Executable()

	log.Println("Started at ", pathToExe)

	if len(os.Args) == 0 {
		log.Fatalln("You must provide a full path to file that will be converted")
	}

	file, err := os.Open(os.Args[1])

	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	if err != nil {
		panic(errors.New(fmt.Sprintf("Could not open file %s", os.Args[1])))
	}

	pureFileName := filepath.Base(os.Args[1])
	pureFileName = pureFileName[:strings.LastIndex(pureFileName, ".")]

	scanner := bufio.NewScanner(file)

	const msiExeFileName string = "Msiexec.exe"
	var i int = 0

	for scanner.Scan() {

		line := scanner.Text()

		if strings.Contains(line, msiExeFileName) {

			log.Printf("Found '%v' line at: %v", msiExeFileName, i+1)

			data := extractData(line)

			log.Println("Detected elements: ", data)

			outputFile, err := os.OpenFile(pureFileName+".out", syscall.O_CREAT|syscall.O_TRUNC, 666)
			if err != nil {
				log.Fatal("There was an exception in opening output file: ", pureFileName+".out")
			}

			var keys []string

			for k := range data {
				keys = append(keys, k)
			}

			sort.Strings(keys)

			for _, key := range keys {
				outputFile.WriteString(fmt.Sprintf("%v=%v\n", key, data[key]))
			}

			break
		}

		i++
	}
}

func extractData(line string) map[string]string {
	elements := make(map[string]string)

	data := strings.Split(line, " ")

	for _, item := range data {
		singleItem := strings.Split(item, "=")

		if len(singleItem) > 1 {
			key := singleItem[0]
			val := singleItem[1]

			val = strings.ReplaceAll(val, "\"", "")

			elements[key] = val
		}
	}

	return elements
}
