package main

import (
	"os"
)

type Page struct {
	Title string
	Body  []byte // "a byte slice" expected by the io libraries we will use

}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}
