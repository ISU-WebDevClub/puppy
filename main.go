package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseFiles("edit.html", "dog.html"))

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
	dog, err := loadDog(name)
	if err != nil {
		http.Redirect(w, r, "/edit/"+name, http.StatusFound)
		return
	}
	renderTemplate(w, "dog", dog)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/edit/"):]
	dog, err := loadDog(name)
	if err != nil {
		dog = &Dog{Name: name}
	}
	renderTemplate(w, "edit", dog)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/save/"):]
	about := r.FormValue("about")
	dog := &Dog{Name: name, About: []byte(about)}
	err := dog.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/dog/"+name, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, dog *Dog) {
	err := templates.ExecuteTemplate(w, tmpl+".html", dog)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/dog/", dogHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
