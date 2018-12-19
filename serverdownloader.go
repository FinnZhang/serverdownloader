package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func home(w http.ResponseWriter, r *http.Request) {
	bs, err := ioutil.ReadFile("./index.html")
	if err != nil {
		log.Println(err)
		return
	}
	_, err = w.Write(bs)
	if err != nil {
		log.Println(err)
		return
	}
}

func download(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if len(strings.Trim(url, " ")) == 0 {
		return
	}
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	//get header
	fileName := resp.Header.Get("Content-Disposition")
	contentType := resp.Header.Get("Content-Type")
	contentLen := resp.Header.Get("Content-Length")
	if fileName == "" {
		fileName = filepath.Base(url)
	}
	//set new header
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", contentLen)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	fileName := header.Filename
	if _, err := os.Stat(".\\static"); os.IsNotExist(err) {
		err := os.Mkdir(".\\static", 0777)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	newFile, err := os.Create("./static/" + fileName)
	defer newFile.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.Copy(newFile, file)
}

func list(w http.ResponseWriter, r *http.Request) {
	html := `<ul>`
	files, err := ioutil.ReadDir("static")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, file := range files {
		html += `<li><a href="/static/` + file.Name() + `">` + file.Name() + `</a></li>`
	}
	html += "</ul>"
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = w.Write([]byte(html))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func clear(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat(".\\static"); !os.IsNotExist(err) {
		err := os.RemoveAll(".\\static")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", home)
	http.HandleFunc("/dl", download)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/list", list)
	http.HandleFunc("/clear", clear)
	http.ListenAndServe(":80", nil)

}
