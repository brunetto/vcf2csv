package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/brunetto/goutils/debug"
	"github.com/brunetto/goutils/readfile"
)

func main() {
	defer debug.TimeMe(time.Now())

	var (
		header         string = "Name,Given Name,Additional Name,Family Name,Yomi Name,Given Name Yomi,Additional Name Yomi,Family Name Yomi,Name Prefix,Name Suffix,Initials,Nickname,Short Name,Maiden Name,Birthday,Gender,Location,Billing Information,Directory Server,Mileage,Occupation,Hobby,Sensitivity,Priority,Subject,Notes,Group Membership,E-mail 1 - Type,E-mail 1 - Value,E-mail 2 - Type,E-mail 2 - Value,IM 1 - Type,IM 1 - Service,IM 1 - Value,Phone 1 - Type,Phone 1 - Value\n"
		fill           string = ",,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,"
		inFiles        []string
		inFileName     string
		line           string
		nReader        *bufio.Reader
		nWriter        *bufio.Writer
		inFile         *os.File
		outFile        *os.File
		err            error
		regNameString  string         = `N:;(.+);;;`
		regPhoneString string         = `TEL;CELL:(\+*\d+)`
		regNameReg     *regexp.Regexp = regexp.MustCompile(regNameString)
		regPhoneReg    *regexp.Regexp = regexp.MustCompile(regPhoneString)
		regNameRes     []string
		regPhoneRes    []string
		data           = map[string]string{}
		name, phone    string
	)

	if inFiles, err = filepath.Glob("*.vcf"); err != nil {
		log.Fatal("Error globbing files in this folder: ", err)
	}

	for _, inFileName = range inFiles {
		// Open infile for reading
		if inFile, err = os.Open(inFileName); err != nil {
			log.Fatal(err)
		}
		defer inFile.Close()
		nReader = bufio.NewReader(inFile)
		// Scan lines
		for {
			if line, err = readfile.Readln(nReader); err != nil {
				break
			}
			if regNameRes = regNameReg.FindStringSubmatch(line); regNameRes != nil {
				name = regNameRes[1]
			} else if regPhoneRes = regPhoneReg.FindStringSubmatch(line); regPhoneRes != nil {
				phone = strings.TrimSpace(regPhoneRes[1])
			}
		}
		data[name] = phone
	}

	// Open file for writing
	outFile, err = os.Create("out.csv")
	defer outFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	nWriter = bufio.NewWriter(outFile)
	defer nWriter.Flush()

	// Write decreased ints to file
	if _, err = nWriter.WriteString(header); err != nil {
		log.Fatalf("Can't write to %v with error %v\n", "out.csv", err)
	}

	for key, value := range data {
		if _, err = nWriter.WriteString(key + "," + key + fill + ",Mobile," + value + "\n"); err != nil {
			log.Fatalf("Can't write to %v with error %v\n", "out.csv", err)
		}
	}

}
