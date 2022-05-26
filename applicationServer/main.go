package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"
)

var (
	serverIP        = flag.String("s", "", "Listen IP")
	serverPort      = flag.String("p", "8080", "Listen Port")
	wordsServerIP   = flag.String("w", "", "Words Server IP")
	wordsServerPort = flag.String("wp", "8081", "Words Server Port")
)

type Word struct {
	Word    string `json:"word"`
	Meaning string `json:"meaning"`
}

func main() {
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/words/{word}", searchWord).Methods("GET")
	r.HandleFunc("/updateWords", updateWordsHandler).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	http.ListenAndServe(*serverIP+":"+*serverPort, r)
}

func searchWord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	word := vars["word"]
	jsonFile, err := os.Open("words/en.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		w.Write([]byte("Something went wrong"))
		return
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var words []Word

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &words)
	for _, y := range words {
		if y.Word == word {
			var word Word
			word.Word = y.Word
			word.Meaning = y.Meaning
			err := json.NewEncoder(w).Encode(word)
			if err != nil {
				w.Write([]byte("Something went wrong"))
				return
			}
			return
		}
	}
	w.Write([]byte("Words Server is Down"))
	return

}
func updateWordsHandler(w http.ResponseWriter, r *http.Request) {
	err := wget("http://" + *wordsServerIP + ":" + *wordsServerPort + "/words.tar.gz")
	if err != nil {
		fmt.Println(err.Error())
		w.Write([]byte("Words Server is Down"))
		return
	}
	cmd := exec.Command("tar", "-xvf", "./words/words.tar.gz", "-C", "./words/")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		w.Write([]byte("Woops something went wrong"))
		return
	}
	err = removeFile("./words/words.tar.gz")
	if err != nil {
		w.Write([]byte("Woops something went wrong"))
		return
	}
	w.Write([]byte("Words Updated"))
}

func wget(fullURLFile string) error {

	// Build fileName from fullPath
	fileURL, err := url.Parse(fullURLFile)
	if err != nil {
		return err
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]
	fmt.Println(fileName)
	fileDir := "./words/"
	// Create blank file
	file, err := os.Create(fileDir + fileName)
	if err != nil {
		return err
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Put content on file
	resp, err := client.Get(fullURLFile)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err

	}
	defer file.Close()
	return nil
}

func removeFile(filename string) error {
	err := os.Remove(filename)
	if err != nil {
		return err
	}
	return nil
}
