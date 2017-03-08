package justhtml

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	templatesDir = "templates"
	pagesDir     = "pages"
	buildDir     = "www"
)

// CreateSite initializes a JustHTML site in the current directory. It returns
// an error if anything goes wrong during the process (which will seriously
// mess things up down the line, so they should be reported back to the user
// and the program stopped).
func CreateSite() error {
	// Create directories
	err := os.Mkdir(templatesDir, 0700)
	if err != nil {
		return fmt.Errorf("Unable to create templates/: %s\n", err)
	}

	err = os.Mkdir(pagesDir, 0700)
	if err != nil {
		return fmt.Errorf("Unable to create pages/: %s\n", err)
	}

	// Create base templates
	hf, err := openFile(templatesDir, "header.tmpl")
	if err != nil {
		return fmt.Errorf("Unable to open header.tmpl: %s\n", err)
	}
	defer hf.Close()
	ff, err := openFile(templatesDir, "footer.tmpl")
	if err != nil {
		return fmt.Errorf("Unable to open footer.tmpl: %s\n", err)
	}
	defer ff.Close()

	h := `{{define "header"}}<!DOCTYPE HTML>
<html>
<head>
	<meta charset="utf-8">

	<title>Title</title>

	<link rel="stylesheet" type="text/css" href="/css/css.css" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</head>
<body>
{{end}}`

	f := `{{define "footer"}}
	{{template "foot" .}}
	{{template "body-end" .}}
{{end}}
{{define "foot"}}



{{end}}
{{define "body-end"}}</body></html>{{end}}`

	if _, err = hf.WriteString(h); err != nil {
		return fmt.Errorf("Unable to write to new file: %s\n", err)
	}

	if _, err = ff.WriteString(f); err != nil {
		return fmt.Errorf("Unable to write to new file: %s\n", err)
	}

	return nil
}

// CreateNewPage creates a file with the given name filled with the normal
// page template.
func CreateNewPage(name string) error {
	f, err := openFile(pagesDir, name+".html")
	if err != nil {
		return fmt.Errorf("Unable to open new file: %s\n", err)
	}
	defer f.Close()

	// Write new basic page
	l := `{{define "` + name + `"}}
{{template "header" .}}



{{template "body-end" .}}
{{end}}`

	if _, err = f.WriteString(l); err != nil {
		return fmt.Errorf("Unable to write to new file: %s\n", err)
	}

	return nil
}

// BuildSite compiles the templates and pages in the current project directory
// into a full static site.
func BuildSite() error {
	start := time.Now()

	// Validate that we're in the correct place with everything we need.
	if _, err := os.Stat(templatesDir); err != nil {
		return fmt.Errorf("FAILED. No templates directory found.")
	}
	if _, err := os.Stat(pagesDir); err != nil {
		return fmt.Errorf("FAILED. No pages directory found.")
	}

	// Create build destination
	fmt.Printf("Creating build directory...")
	err := os.Mkdir(buildDir, 0700)
	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			fmt.Printf("SKIP: already exists\n")
		} else {
			return fmt.Errorf("Unable to create %s/: %s\n", buildDir, err)
		}
	} else {
		fmt.Printf("Done\n")
	}

	// Generate pages
	filepath.Walk(pagesDir, func(path string, i os.FileInfo, err error) error {
		if !i.IsDir() && !strings.HasPrefix(i.Name(), ".") {
			fmt.Printf("Opening page %s...", i.Name())
			f, err := openFile(buildDir, i.Name())
			if err != nil {
				err = fmt.Errorf("SKIP: unable to open file in build directory (%s).", buildDir)
				fmt.Printf("%s\n", err)
				return err
			}
			defer f.Close()

			w := bufio.NewWriter(f)
			template, err := template.ParseFiles(filepath.Join(pagesDir, i.Name()), filepath.Join(templatesDir, "header.tmpl"), filepath.Join(templatesDir, "footer.tmpl"))
			if err != nil {
				fmt.Printf("SKIP: %s\n", err)
				return err
			}
			template.ExecuteTemplate(w, i.Name()[:strings.LastIndex(i.Name(), ".")], nil)
			w.Flush()
			fmt.Println("Done")
		}

		return nil
	})

	fmt.Printf("Build finished in %s\n", time.Since(start))
	return nil
}

func openFile(d, n string) (*os.File, error) {
	return os.OpenFile(filepath.Join(d, n), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
}
