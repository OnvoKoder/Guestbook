package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type Guestbook struct {
	Count int
	List  []string
}

type Comment struct {
	Author  string
	Date    string
	Message string
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetComments(fileName string) ([]string, error) {
	comment, err := ReadJson(fileName)
	if err != nil {
		check(err)
	}
	array := make([]string, 0)
	for _, com := range comment {
		array = append(array, com.Author+" "+com.Date+" "+com.Message)
	}
	return array, nil
}

func ReadJson(fileName string) ([]Comment, error) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		check(err)
	}
	var comment []Comment
	json.Unmarshal([]byte(file), &comment)
	return comment, nil
}

func addCommentHandler(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("addComment.html")
	check(err)
	err = html.Execute(writer, nil)
	check(err)
}

func createtHandler(writer http.ResponseWriter, request *http.Request) {
	message := request.FormValue("comment")
	author := request.FormValue("author")
	currentTime := time.Now()
	if author == "" {
		author = "unknown"
	}
	comment, err := ReadJson("test.json")
	if err != nil {
		check(err)
	}
	comment = append(comment, Comment{Author: author, Date: currentTime.Format("2006-01-02 15:04:05"), Message: message})
	data, _ := json.Marshal(comment)
	_ = os.WriteFile("test.json", data, 0644)
	http.Redirect(writer, request, "/guestbook", http.StatusFound)
}

func viewHandler(writer http.ResponseWriter, request *http.Request) {
	signature, err := GetComments("test.json")
	check(err)
	html, err := template.ParseFiles("index.html")
	check(err)
	guestbook := Guestbook{len(signature), signature}
	err = html.Execute(writer, guestbook)
	check(err)
}

func main() {
	fmt.Println("Server run...")
	http.HandleFunc("/guestbook", viewHandler)
	http.HandleFunc("/guestbook/addComment", addCommentHandler)
	http.HandleFunc("/guestbook/create", createtHandler)
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
