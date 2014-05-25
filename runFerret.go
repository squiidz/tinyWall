package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Post struct {
	Title    string
	Content  template.HTML
	Date     time.Time
	Comments []Comment
}

type Comment struct {
	Username string
	Content  string
	Date     time.Time
	Like     int
}

func ShitAppend(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func NewPost(title string, content string, date time.Time) Post {
	post := Post{
		Title:   title,
		Content: template.HTML(content),
		Date:    date}
	return post
}

func LoadFile(file string) []byte {
	doc, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func (post *Post) SaveFile() {
	err := ioutil.WriteFile("post/"+post.Title+".txt", []byte(post.Content), 0777)
	if err != nil {
		fmt.Println("[*] Post cannot be Saved !")
	}
}

func FindPost() []Post {
	posts := []Post{}
	wd, err := os.Getwd()
	ShitAppend(err)
	files, err := ioutil.ReadDir(wd + "/post/") // Find Post In Working Directory + /post

	for loop, doc := range files {

		if doc.IsDir() != true && strings.Contains(doc.Name(), ".txt") == true {
			access := "post/" + doc.Name()
			cont, err := ioutil.ReadFile(access)
			if err != nil {
				fmt.Println(err)
			}
			posts = append(posts, NewPost(files[loop].Name(), string(cont), files[loop].ModTime()))
		}

	}
	return posts
}

func main() {
	http.HandleFunc("/", Handler)
	http.HandleFunc("/edit", EditHandler)
	http.HandleFunc("/add", PostHandler)

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("fonts"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Handler(rw http.ResponseWriter, req *http.Request) {
	log.Println("[+] " + req.RemoteAddr + " is connected [+]")
	if req.URL.Path != "/" { // Check if the request is for the root
		http.NotFound(rw, req)
		return
	}
	posts := FindPost()
	temp, _ := template.ParseFiles("template/index.html")
	for loop, _ := range posts {
		filename := strings.Split(posts[loop].Title, ".txt")
		posts[loop].Title = filename[0]
		posts[loop].Content = template.HTML(string(posts[loop].Content))
	}
	temp.Execute(rw, posts)
}

func EditHandler(rw http.ResponseWriter, req *http.Request) {
	file := LoadFile("template/edit.html")
	rw.Write(file)
}

func PostHandler(rw http.ResponseWriter, req *http.Request) {
	newPost := Post{
		Title:   req.FormValue("title"),
		Content: template.HTML(req.FormValue("content")),
		Date:    time.Now()}
	newPost.SaveFile()
	http.Redirect(rw, req, "/", http.StatusFound)
	log.Println("[+] " + req.RemoteAddr + " create a post " + newPost.Title + " [+]")
}
