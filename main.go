package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	excludeFileTypes []string
	excludeDirs      []string
	linesEachFile    map[string]int
	scannerBuffer    *int
	totalSum         = 0
)

func readFile(path string) {

	file, err := os.Open(path)
	defer file.Close() //close when done
	if err != nil {
		log.Fatalf(err.Error())
	}
	scanner := bufio.NewScanner(file)
	buf := make([]byte, *scannerBuffer)
	scanner.Buffer(buf, *scannerBuffer)
	scanner.Split(bufio.ScanLines)
	linesEachFile[path] = 0
	for scanner.Scan() {
		if len(scanner.Text()) > 0 { //skip if it's just an empty line
			linesEachFile[path]++
			totalSum++
		}
	}
}

func iterateOverDir(path string) {
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		skip := false
		if err != nil {
			log.Fatalf(err.Error())
		}

		if !f.IsDir() {
			for _, e := range excludeDirs {
				if strings.Contains(path, e) {
					skip = true
				}
			}
			if !skip {
				extension := strings.Split(f.Name(), ".")
				if len(extension) > 1 {
					for _, e := range excludeFileTypes {
						if e == extension[1] {
							skip = true
							break
						}
					}
				}

				if !skip {
					readFile(path)
				}
			}
			skip = false

		}
		return nil
	})
}

func printResult() {
	if len(linesEachFile) == 0 || totalSum == 0 {
		fmt.Println("No files were read!")
		return
	}
	for k, e := range linesEachFile {
		fmt.Printf("File %s contained %d lines\n", k, e)
	}
	fmt.Printf("This project has a total lenght of %d lines\n", totalSum)
}

func main() {
	dirFlag := flag.String(
		"dir",
		"None",
		"Set the directory to iterate over",
	)
	excludeFileFlag := flag.String(
		"excludeFiles",
		"None",
		"Add file extensions to exclude.\n-> Split with ;\n-> no . needed",
	)
	excludeDirsFlag := flag.String(
		"excludeDirs",
		"None",
		"Add directories that shall be excluded.\n-> Split with ;",
	)
	scannerBuffer = flag.Int(
		"scannerBuffer",
		64000,
		"Adjust the size of the scanner's buffer, when reading a file (in lines).",
	)
	flag.Parse()
	if *dirFlag == "None" {
		log.Fatalf("No directoy was given!")
	}
	excludeFileTypes = strings.Split(*excludeFileFlag, ";")
	excludeDirs = strings.Split(*excludeDirsFlag, ";")
	linesEachFile = make(map[string]int)
	iterateOverDir(*dirFlag)
	printResult()
}
