package main

import (
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type listItem struct {
	FileName, URL string
}

type listPageData struct {
	Items []listItem
}

//Init log
func Init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func home(w http.ResponseWriter, r *http.Request) {
	bs, err := ioutil.ReadFile("./index.html")
	if err != nil {
		logHTTPErr(w, err)
		return
	}
	_, err = w.Write(bs)
	if err != nil {
		logHTTPErr(w, err)
	}
}

func download(w http.ResponseWriter, r *http.Request) {
	url := strings.TrimSpace(r.FormValue("url"))
	if len(url) == 0 {
		err := errors.New("Download link is invalid. ")
		logHTTPErr(w, err)
		return
	}
	resp, err := http.Get(url)
	if err != nil {
		logHTTPErr(w, err)
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
		logHTTPErr(w, err)
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		logHTTPErr(w, err)
		return
	}
	defer file.Close()
	fileName := header.Filename
	err = createStaitcDirIfNotExist()
	if err != nil {
		logHTTPErr(w, err)
		return
	}

	newFile, err := os.Create("./static/" + fileName)
	defer newFile.Close()
	if err != nil {
		logHTTPErr(w, err)
		return
	}
	io.Copy(newFile, file)
	http.Redirect(w, r, "/list", http.StatusTemporaryRedirect)
}

func list(w http.ResponseWriter, r *http.Request) {
	err := createStaitcDirIfNotExist()
	if err != nil {
		logHTTPErr(w, err)
		return
	}
	files, err := ioutil.ReadDir("./static")
	if err != nil {
		logHTTPErr(w, err)
		return
	}
	tmpl := template.Must(template.ParseFiles("./tmpl/list.html"))
	var data listPageData
	for _, file := range files {
		fileName := file.Name()
		item := listItem{FileName: fileName, URL: "/static/" + fileName}
		data.Items = append(data.Items, item)
	}

	if len(data.Items) == 0 {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = tmpl.Execute(w, data)
		if err != nil {
			logHTTPErr(w, err)
		}
	}
}

func clear(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat("./static"); !os.IsNotExist(err) {
		err := os.RemoveAll("./static")
		if err != nil {
			logHTTPErr(w, err)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func createStaitcDirIfNotExist() error {
	if _, err := os.Stat("./static"); os.IsNotExist(err) {
		err = os.Mkdir("./static", 0777)
		return err
	}
	return nil
}

func logHTTPErr(w http.ResponseWriter, err error) {
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	//write log in a file named log.txt
	file, err := os.OpenFile("./log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)
	log.Println("****************************start*************************")
	log.Println(time.Now().String())

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", home)
	http.HandleFunc("/dl", download)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/list", list)
	http.HandleFunc("/clear", clear)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
