package main

import (
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseFiles("edit.html", "dog.html"))
var validPath = regexp.MustCompile("^/(edit|save|dog)/([a-zA-Z0-9]+)$")

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

// makeHandler returns a function of type http.HandlerFunc. This function
// validates the name and then passes the name to the function provided as an
// argument to makeHandler. This allows us to remove the duplicate function calls
// to validate the name in each handler.
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

// dogHandler runs when we hit the "/dog/<name>" endpoint on the HTTP server.
func dogHandler(w http.ResponseWriter, r *http.Request, name string) {
	dog, err := loadDog(name)
	if err != nil {
		http.Redirect(w, r, "/edit/"+name, http.StatusFound)
		return
	}
	renderTemplate(w, "dog", dog)
}

// editHandler runs when we hit the "/edit/<name>" endpoint on the HTTP server.
func editHandler(w http.ResponseWriter, r *http.Request, name string) {
	dog, err := loadDog(name)
	if err != nil {
		dog = &Dog{Name: name}
	}
	renderTemplate(w, "edit", dog)
}

// saveHandler takes information from the HTML form and creates a new dog.
func saveHandler(w http.ResponseWriter, r *http.Request, name string) {
	about := r.FormValue("about")
	dog := &Dog{Name: name, About: []byte(about)}
	err := dog.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/dog/"+name, http.StatusFound)
}

// renderTemplate is a helper function for executing a given template with the
// given dog.
func renderTemplate(w http.ResponseWriter, tmpl string, dog *Dog) {
	err := templates.ExecuteTemplate(w, tmpl+".html", dog)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// getName validates the dog name in the URL. This is important because this user
// input is written to the file system.
func getName(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid dog name")
	}
	return m[2], nil // m[2] is the dog's name
}

func main() {
	http.HandleFunc("/dog/", makeHandler(dogHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
