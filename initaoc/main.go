package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	year := flag.Int("year", time.Now().Year(), "the year")
	day := flag.Int("day", time.Now().Day(), "day number")
	sessionFilename := flag.String("session", "session.txt", "the filename holding the AoC session key")
	boilerplateFilename := flag.String("b", "boilerplate.txt", "boilerplate filename")
	pbaseUrl := flag.String("url", "https://adventofcode.com", "aoc url")
	pinput := flag.String("input", "input.txt", "input filename")
	pmain := flag.String("main", "main.go", "main go filename")
	flag.Parse()

	session, err := readFile(*sessionFilename)
	if err != nil {
		log.Fatal(err)
	}

	boilerplate, err := readFile(*boilerplateFilename)
	if err != nil {
		log.Fatal(err)
	}

	err = runInit(*year, *day, session, boilerplate, *pbaseUrl, *pinput, *pmain)
	if err != nil {
		log.Fatal(fmt.Sprintf("couldn't init aoc with the params day: %d year %d err: %s", *day, *year, err.Error()))
	}
}

func runInit(year, day int, session, boilerplate, baseUrl, inputFilename, mainFilename string) error {
	syear, sday := strconv.Itoa(year), strconv.Itoa(day)
	pathDay := "0" + sday
	pathDay = pathDay[len(pathDay)-2:]
	inputPath := path.Join(syear, "day", sday, "input")

	fullUrl := strings.Join([]string{baseUrl, inputPath}, "/")
	input, err := getInput(fullUrl, session)
	if err != nil {
		return fmt.Errorf("can't get input at %s %w", fullUrl, err)
	}

	dirPath := filepath.Join(syear, pathDay)
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}

	err = initFile(dirPath, mainFilename, boilerplate, true)
	if err != nil {
		log.Println("main file error", err)
	} else {
		log.Printf("init'd file %s\\%s", dirPath, mainFilename)
	}
	err = initFile(dirPath, inputFilename, input, false)
	if err != nil {
		log.Println("input file error", err)
	} else {
		log.Printf("init'd file %s\\%s", dirPath, inputFilename)
	}

	return nil
}

func initFile(dir, filename, contents string, failIfExists bool) error {
	fpath := filepath.Join(dir, filename)
	_, err := os.Stat(fpath)

	if failIfExists && (os.IsExist(err) || err == nil) {
		return fmt.Errorf("i won't overwrite a file that already exists. %s", fpath)
	}

	f, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(contents)

	return err
}

func getInput(url, session string) (string, error) {
	c := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("cookie", fmt.Sprintf("session=%s", session))

	res, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("couldn't get file contents at url: %s  status: %s", url, res.Status)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func readFile(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	b, err := io.ReadAll(f)

	return string(b), err
}
