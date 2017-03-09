package main

import (
	"fmt"
	"github.com/thebaer/justhtml"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		showHelp()
	}
	cmd := os.Args[1]

	switch cmd {
	case "init":
		err := justhtml.CreateSite()
		handleErr(err)
		break
	case "new":
		if !commandHasArgs(1) {
			showUsage("new [page-name]")
		}
		err := justhtml.CreateNewPage(os.Args[2])
		handleErr(err)
		break
	case "build":
		err := justhtml.BuildSite()
		handleErr(err)
		break
	default:
		showHelp()
	}
}

func handleErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func showUsage(u string) {
	fmt.Printf("usage: jhtml %s\n", u)
	os.Exit(1)
}

func commandHasArgs(n int) bool {
	return len(os.Args) >= n+2
}

func showHelp() {
	fmt.Println(`usage: jhtml <command> [<args>]

Commands:

    init        Create project directories in current location
    new [name]  Create a new page with the given name, e.g. index
    build       Generate static site from current project directory
`)

	os.Exit(1)
}
