package justhtml

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	templatesDir = "templates"
	pagesDir     = "pages"
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

func openFile(d, n string) (*os.File, error) {
	return os.OpenFile(filepath.Join(d, n), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
}
