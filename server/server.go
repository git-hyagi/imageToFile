package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/otiai10/gosseract/v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

const outputMetadata = `
  <form method="POST" enctype="multipart/form-data">
  <input type="file" name="imageToText">
  <input type="submit">
  </form>
  <br>`

func main() {
	http.HandleFunc("/", webPage)
	http.ListenAndServe(":8080", nil)
}

func webPage(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {

		file, handler, err := req.FormFile("imageToText")
		if err != nil {
			fmt.Println("Error (1)!" + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		fmt.Printf("Received file: %v\n", handler.Filename)
		fmt.Printf("Size: %v\n", handler.Size)
		fmt.Printf("Header: %v\n", handler.Header)

		newFile, err := os.Create("/tmp/tmp.png")
		if err != nil {
			fmt.Println("Error (2)!" + err.Error())
			return
		}
		fileToBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println("Error (3)!" + err.Error())
			return
		}
		_, err = newFile.Write(fileToBytes)
		if err != nil {
			fmt.Println("Error (4)!" + err.Error())
			return
		}
	}

	// convert image to text
	makeText("/tmp/tmp.png", "/tmp/tmp.out")

	// open the converted file
	textedFile, err := os.Open("/tmp/tmp.out")
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(textedFile)

	b := new(bytes.Buffer)
	fmt.Fprint(b, "<table style=\"width:50%\">")
	fmt.Fprint(b, "<tr><th>DATE</th><th>STORE</th><th>VALUE</th><th>KIND</th></tr>")

	// parse the file according to regex
	regexCharge := regexp.MustCompile(`^.*?(?P<date>\d{2}/\d{2})\s+(?P<store>.*)\s+(?P<value>[-]?\d+(\.\d{3})?,\d{2})$`)
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
			fmt.Fprint(b, "<tr>")
			fmt.Fprintf(b, "<td>%v</td> <td>%v</td> <td>%v</td> <td>%v</td>", groups["date"], groups["store"], groups["value"], scanner.Text())
			fmt.Fprint(b, "</tr>")
		}
	}
	fmt.Fprint(b, "</table>")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, outputMetadata+b.String())
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
