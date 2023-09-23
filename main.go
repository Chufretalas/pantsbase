package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/Chufretalas/pantsbase/db"
	"github.com/Chufretalas/pantsbase/routes"
	"golang.org/x/exp/slices"
)

func main() {
	db.ConnectDB()
	defer db.DB.Close()

	routes.LoadRoutes()
	fmt.Println("Listening on http://localhost:8000")

	if slices.Contains(os.Args, "--open") {
		openBrowser()
	}

	http.ListenAndServe(":8000", routes.Router)
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
