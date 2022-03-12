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

// dict[aa] de ex si returneaza dictul cu cuvintele care incep cu aa
//var dictionar map[string]cuvinte
var dictionar = make(map[string]cuvinte)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

//generate the dict here
func init() {
	fmt.Println("Starting server...")

	files, err := ioutil.ReadDir("./data/")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		dictionar[file.Name()] = makeCuvinteDict("./data/" + file.Name())
	}
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

	fmt.Println("Server started!")
	http.ListenAndServe(":1337", nil)
}

func ServeWord(w http.ResponseWriter, r *http.Request) {
	Cuvinte := dictionar[r.URL.Query()["starting"][0]]

	randomIndex := random.Intn(len(Cuvinte))

	word := "N/A"
	i := 0
	for k, _ := range Cuvinte {
		i += 1
		if i == randomIndex {
			word = k
			break
		}
	}

	fmt.Fprintf(w, word)
}

func IsValidWord(w http.ResponseWriter, r *http.Request) {
	cuvant := r.URL.Query()["word"][0]
	//cuvant := dictionar[r.URL.Query()["word"][0]]

	var sb strings.Builder
	sb.WriteByte(cuvant[0])
	sb.WriteByte(cuvant[1])

	for k, _ := range dictionar[sb.String()] {
		if k == cuvant {
			fmt.Fprintf(w, "true")
			return
		}
	}
	fmt.Fprintf(w, "false")
}
