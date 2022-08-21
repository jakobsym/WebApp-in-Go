/* To run ->  http://localhost:8080/view/ANewPage
'ANewPage' can be replaced with whatever file name you may desire. */
package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
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

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

/* The function template.Must is a convenience wrapper that panics when passed a non-nil error value,
and otherwise returns the *Template unaltered. A panic is appropriate here; if the templates can't be loaded the only sensible thing to do is exit the program.*/
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	/* Closure, because this func encloses values defined outside of it. In our case,
	'fn' is enclosed by the closure. 'fn' will be our (edit, save, or view) */
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler)) // VIEW
	http.HandleFunc("/edit/", makeHandler(editHandler)) // EDIT
	http.HandleFunc("/save/", makeHandler(saveHandler)) // SAVE
	log.Fatal((http.ListenAndServe(":8080", nil)))

}
