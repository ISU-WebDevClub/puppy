package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Dog is a data structure to represent a dog.
type Dog struct {
	Name  string
	About []byte
	// About is a byte slice rather than a string so that io libraries can use it.
}

// save is a method that, when called on a Dog, creates a new file with the dog's
// name as the file name and the about as the contents of the file.
func (p *Dog) save() error {
	filename := p.Name + ".dog"
	return ioutil.WriteFile(filename, p.About, 0600)
}

// loadDog reads information from the file with the given name and returns a
// pointer to a new Dog struct.
func loadDog(name string) (*Dog, error) {
	filename := name + ".dog"
	about, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Dog{Name: name, About: about}, nil
}

// dogHandler runs when we hit the "/dog/" endpoint on the HTTP server.
func dogHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/dog/"):]
	dog, _ := loadDog(name)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", dog.Name, dog.About)
}

func main() {
	http.HandleFunc("/dog/", dogHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
