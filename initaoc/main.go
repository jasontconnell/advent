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
	"regexp"
	"strconv"
	"strings"
	"time"
)

type file struct {
	filename string
	contents string
}

func main() {
	year := flag.Int("y", time.Now().Year(), "the year")
	day := flag.Int("d", time.Now().Day(), "day number")
	sessionFilename := flag.String("session", "session.txt", "the filename holding the AoC session key")
	uaFilename := flag.String("useragent", "useragent.txt", "the filename holding the user agent (new req as of 2022)")
	boilerplateFolder := flag.String("b", "./boilerplate/", "boilerplate folder")
	pbaseUrl := flag.String("url", "https://adventofcode.com", "aoc url")
	pinput := flag.String("input", "input.txt", "input filename")
	flag.Parse()

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(fmt.Errorf("getwd failed for some reason %w", err))
	}

	y := *year
	// check if we're in a year folder
	_, folder := filepath.Split(cwd)
	reg := regexp.MustCompile("([0-9]{4})")
	m := reg.FindStringSubmatch(folder)
	createYearDir := true
	if len(m) == 2 {
		y, _ = strconv.Atoi(folder)
		createYearDir = false
	}

	session, err := readFile(*sessionFilename)
	if err != nil {
		log.Fatal(err)
	}

	useragent, err := readFile(*uaFilename)
	if err != nil {
		log.Fatal(err)
	}

	err = runInit(y, *day, createYearDir, session, useragent, *boilerplateFolder, *pbaseUrl, *pinput)
	if err != nil {
		log.Fatalf("couldn't init aoc with the params day: %d year %d err: %s", *day, *year, err.Error())
	}
}

func runInit(year, day int, createYearDir bool, session, useragent, boilerplate, baseUrl, inputFilename string) error {
	syear, sday := strconv.Itoa(year), strconv.Itoa(day)
	pathDay := "0" + sday
	pathDay = pathDay[len(pathDay)-2:]
	inputPath := path.Join(syear, "day", sday, "input")

	fullUrl := strings.Join([]string{baseUrl, inputPath}, "/")
	input, err := getInput(fullUrl, session, useragent)
	if err != nil {
		return fmt.Errorf("can't get input at %s %w", fullUrl, err)
	}

	dirPath := filepath.Join(syear, pathDay)
	if !createYearDir {
		dirPath = pathDay
	}

	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}

	files, err := readFolder(boilerplate)
	if err != nil {
		return err
	}

	files = append(files, file{filename: inputFilename, contents: input})

	for _, f := range files {
		contents := f.contents
		contents = strings.ReplaceAll(contents, "{year}", syear)
		contents = strings.ReplaceAll(contents, "{day}", pathDay)

		success, err := initFile(dirPath, f.filename, contents, true)
		if err != nil {
			log.Printf("error creating file %s\\%s, %s", dirPath, f.filename, err.Error())
		} else if success {
			log.Printf("init'd file %s\\%s", dirPath, f.filename)
		} else {
			log.Printf("skipped file %s\\%s", dirPath, f.filename)
		}
	}

	return nil
}

func initFile(dir, filename, contents string, skipIfExists bool) (bool, error) {
	fpath := filepath.Join(dir, filename)
	_, err := os.Stat(fpath)

	if skipIfExists && (os.IsExist(err) || err == nil) {
		return false, nil
	}

	f, err := os.Create(fpath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.WriteString(contents)

	return err == nil, err
}

func getInput(url, session, useragent string) (string, error) {
	c := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("cookie", fmt.Sprintf("session=%s", session))
	req.Header.Add("User-Agent", useragent)

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

func readFolder(folder string) ([]file, error) {
	files := []file{}
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		_, fn := filepath.Split(path)

		f := file{filename: fn}

		b, rerr := os.ReadFile(path)
		if rerr != nil {
			return err
		}
		f.contents = string(b)
		files = append(files, f)
		return nil
	})
	return files, err
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
