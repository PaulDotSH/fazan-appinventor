package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	_ "embed"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"
)

type cuvinte map[string]byte

var dictionar = make(map[string]cuvinte)

var StartingLetters string

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// The word lists live compressed in data.tar.gz (one file per starting
// bigram). Embedding the archive keeps the repository small; it is
// decompressed in memory at startup instead of shipping an uncompressed
// data/ directory.
//
//go:embed data.tar.gz
var dataArchive []byte

// generate the dict here
func init() {
	fmt.Println("Starting server...")

	contents, err := readDataArchive(dataArchive)
	if err != nil {
		log.Fatal(err)
	}

	names := make([]string, 0, len(contents))
	for name := range contents {
		names = append(names, name)
	}
	sort.Strings(names)

	var builder strings.Builder
	for _, name := range names {
		builder.WriteString(name)
		builder.WriteRune(' ')
		dictionar[name] = makeCuvinteDict(bytes.NewReader(contents[name]))
	}
	StartingLetters = builder.String()
}

func readDataArchive(raw []byte) (map[string][]byte, error) {
	gz, err := gzip.NewReader(bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	files := make(map[string][]byte)
	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if !header.FileInfo().Mode().IsRegular() {
			continue
		}
		data, err := io.ReadAll(tr)
		if err != nil {
			return nil, err
		}
		files[header.Name] = data
	}

	return files, nil
}

func makeCuvinteDict(r io.Reader) cuvinte {
	dict := make(cuvinte)

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		//for line in file dict[line]=1
		dict[scanner.Text()] = 1
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return dict
}

func main() {
	http.HandleFunc("/GetWordStartingWith", ServeWord)
	http.HandleFunc("/IsValidWord", IsValidWord)
	http.HandleFunc("/GetStarting", ListStartingCombinations)

	fmt.Println("Server started!")
	http.ListenAndServe(":1337", nil)
}

func ServeWord(w http.ResponseWriter, r *http.Request) {
	cuvant := r.URL.Query()["start"][0]
	fmt.Println(r.URL.Query()["start"])
	length := len(cuvant)

	if length < 2 {
		return
	}

	var sb strings.Builder

	if length == 2 {
		sb.WriteByte(cuvant[0])
		sb.WriteByte(cuvant[1])
	} else {
		sb.WriteByte(cuvant[length-2])
		sb.WriteByte(cuvant[length-1])
	}

	Cuvinte := dictionar[sb.String()]
	if len(Cuvinte) == 0 {
		fmt.Fprintf(w, "N/A")
		fmt.Println("N/A")
		return
	}
	randomIndex := random.Intn(len(Cuvinte))

	word := "N/A"
	i := 0
	for k, _ := range Cuvinte {
		if i == randomIndex {
			word = k
			break
		}
		i += 1
	}
	fmt.Println(word)

	fmt.Fprintf(w, word)
}

func IsValidWord(w http.ResponseWriter, r *http.Request) {
	cuvinte := r.URL.Query()["word"]
	if len(cuvinte) < 2 {
		return
	}

	var sb strings.Builder

	length := len(cuvinte[1])
	if length < 2 {
		return
	}

	sb.WriteByte(cuvinte[1][length-2])
	sb.WriteByte(cuvinte[1][length-1])

	if !strings.HasPrefix(cuvinte[0], sb.String()) {
		fmt.Fprintf(w, "false")
		return
	}

	sb.Reset()
	sb.WriteByte(cuvinte[0][0])
	sb.WriteByte(cuvinte[0][1])

	tmp := dictionar[sb.String()]

	b := tmp[cuvinte[0]]
	if b == 1 {
		fmt.Fprintf(w, "true")
		return
	}

	fmt.Fprintf(w, "false")
}

func ListStartingCombinations(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, StartingLetters)
}
