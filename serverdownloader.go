package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func download(c echo.Context) error {
	name := c.FormValue("name")
	wgetCmd := exec.Command("wget", "-P", "./public", name)
	wgetCmd.Run()

	return c.HTML(http.StatusOK, lsCmd())
}

func upload(c echo.Context) error {

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(path.Join("./public", file.Filename))
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return err
	}
	return c.HTML(http.StatusOK, lsCmd())
}

func list(c echo.Context) error {
	html := lsCmd()
	if html == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}
	return c.HTML(http.StatusOK, lsCmd())
}

func lsCmd() string {

	if _, err := os.Stat("./public"); os.IsNotExist(err) {
		os.Mkdir("./public", 0644)
	}

	lsCmd := exec.Command("ls", "public")
	lsOut, err := lsCmd.Output()
	if err != nil {
		panic(err)
	}
	files := strings.Split(string(lsOut), "\n")
	var html string
	for _, file := range files {
		if file == "" {
			continue
		}
		a := `<a href="/static/` + file + `">` + file + `</a>` + `<br>`
		html = html + a
	}
	return html
}

func clear(c echo.Context) error {
	err := os.RemoveAll("./public")
	if err != nil {
		return c.HTML(http.StatusOK, err.Error())
	}
	err = os.Mkdir("./public", 0777)
	if err != nil {
		return c.HTML(http.StatusOK, err.Error())
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.File("index.html")
	})
	e.GET("/clear", clear)
	e.GET("/list", list)
	e.POST("/download", download)
	e.POST("/upload", upload)
	e.Static("/static", "./public")

	e.Logger.Fatal(e.Start(":80"))

}
