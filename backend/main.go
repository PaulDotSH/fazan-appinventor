package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type cuvinte map[string]byte

var dictionar = make(map[string]cuvinte)

var StartingLetters string

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

//generate the dict here
func init() {
	fmt.Println("Starting server...")

	files, err := ioutil.ReadDir("./data/")
	if err != nil {
		log.Fatal(err)
	}

	var builder strings.Builder
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		builder.WriteString(file.Name())
		builder.WriteRune(' ')
		dictionar[file.Name()] = makeCuvinteDict("./data/" + file.Name())
	}
	StartingLetters = builder.String()
}

func makeCuvinteDict(fileName string) cuvinte {
	dict := make(cuvinte)

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

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
