/* Defining our data structure */
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

/* 'Page' will describe how page data will be stored in memory */
type Page struct {
	Title string
	Body  []byte // "a byte slice" expected by the io libraries we will use

}

/* This is a method named 'save' that takes as its receiver p, a pointer to Page . It takes no parameters, and returns a value of type error */
/* This method above will save our structs Body to a .txt file with name as 'Title' */
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600) //using 0600 as file should be created with read-write permissions for user.
	//return nil if all works
}

/* func oldLoadPage(title string) *Page {
	filename := title + ".txt"
	body, _ := os.ReadFile(filename) //(_) symbol is used to throw away the error return value as os.ReadFile returns an error
	return &Page{Title: title, Body: body}
} */

/* Notice how we are using the blank identifier to ignore the error produced by os.ReadFile.
We are going to create an improved loadPage function as this error can be useful. */

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, I love %s!", r.URL.Path[1:])
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal((http.ListenAndServe(":8080", nil)))
	/* p1 := &Page{Title: "TestPage", Body: []byte("Hello, TestPage!")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body)) */
}
