package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
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

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/dl", download)
	http.ListenAndServe(":80", nil)

}
