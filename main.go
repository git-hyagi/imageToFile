package main

import (
	"bufio"
	"fmt"
	"github.com/otiai10/gosseract/v2"
	"os"
	"regexp"
)

const outputTextFile = "toText.out"

func main() {
	var fileName string

	// if file name not passed as an argument, use test.png as default
	if len(os.Args) < 2 {
		fileName = "test.png"
	} else {
		fileName = os.Args[1]
	}

	// convert image to text
	makeText(fileName, outputTextFile)

	// open the converted file
	file, err := os.Open(outputTextFile)
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(file)

	// parse the file according to regex
	regexCharge := regexp.MustCompile(`^.*?(?P<date>\d{2}/\d{2})\s+(?P<store>.*)\s+(?P<value>\d+,\d{2})$`)
	for scanner.Scan() {
		groups := map[string]string{}
		matches := regexCharge.FindStringSubmatch(scanner.Text())

		// put the regex matches into the "groups" hashmap (just to easy the access)
		for i, name := range matches {
			if i == 0 { // skip the first occurrence as it contains the leftmost match of the regular expression in string
				continue
			}
			groups[regexCharge.SubexpNames()[i]] = name
		}

		if len(groups) > 0 {
			// the next line always have information of what kind of ... was this credit
			scanner.Scan()
			fmt.Printf("date: %v\tstore: %v\tvalue: %v\tkind: %v\n", groups["date"], groups["store"], groups["value"], scanner.Text())
		}

	}
}

// Convert the png file to text using gosseract
func makeText(imageFile, outputTextFile string) error {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage(imageFile)

	fmt.Println(`Converting image to file ...`)
	text, err := client.Text()
	if err != nil {
		return err
	}

	newFile, err := os.Create(outputTextFile)
	if err != nil {
		return err
	}
	_, err = newFile.WriteString(text)
	if err != nil {
		return err
	}

	return nil
}
