package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"slices"

	"github.com/Chufretalas/pantsbase/controllers"
	"github.com/Chufretalas/pantsbase/db"
	"github.com/Chufretalas/pantsbase/routes"
)

//go:embed all:static
var staticFS embed.FS

//go:embed all:templates
var tempFS embed.FS

func init() {

	tempFS, err := fs.Sub(tempFS, "templates")

	if err != nil {
		log.Fatal(err)
	}

	controllers.Temps = template.Must(template.ParseFS(tempFS, "*.html"))

	staticFS, err := fs.Sub(staticFS, "static")

	if err != nil {
		log.Fatal(err)
	}

	routes.StaticFS = staticFS

	db.SchemaRegex = regexp.MustCompile(`(\"[A-z\s\d]+\"\s[A-z]+)`)
}

//TODO: WRITE TESTS!!!!!!!!!

func main() {
	db.ConnectDB("./data.db")
	defer db.DB.Close()

	router := http.NewServeMux()
	routes.LoadRoutes(router)
	fmt.Println("Listening on http://localhost:8000")

	if slices.Contains(os.Args, "--open") {
		openBrowser()
	}

	http.ListenAndServe(":8000", router)
}

func openBrowser() {
	// thanks to: https://gist.github.com/sevkin/9798d67b2cb9d07cb05f89f14ba682f8
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, "http://localhost:8000/")
	exec.Command(cmd, args...).Start()
}
