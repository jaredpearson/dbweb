package web

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
)

// determineTemplateDir determines the directory containing the template. This
// will use TEMPLATE_PATH environment or a default path the "templates" directory
// within the current working directory.
func determineTemplateDir() string {
	directory := os.Getenv("TEMPLATE_PATH")
	if len(directory) == 0 {
		d, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		directory = path.Join(d, "templates")
		log.Printf("TEMPLATE_PATH not defined. Using default: %s", directory)
	}
	return directory
}

// GetTemplatePath gets the path of a template by the specific name. Usually this
// is "templatePath + name + .html".
func GetTemplatePath(name string) string {
	return path.Join(templateDir, name+".html")
}

type MainLayoutData interface {
	PageTitle() string
	UserInfo() UserInfo
}
type mainLayoutData struct {
	pageTitle string
	userInfo  UserInfo
}

func (main mainLayoutData) PageTitle() string {
	return main.pageTitle
}
func (main mainLayoutData) UserInfo() UserInfo {
	return main.userInfo
}

func NewMainLayoutData(pageTitle string, userInfo UserInfo) MainLayoutData {
	return mainLayoutData{
		pageTitle: pageTitle,
		userInfo:  userInfo,
	}
}

func ShowTemplateInMainLayout(
	w http.ResponseWriter,
	r *http.Request,
	templateName string,
	data MainLayoutData) {
	ShowTemplateInLayout(w, r, "main", templateName, data)
}

func ShowTemplateInLayout(
	w http.ResponseWriter,
	r *http.Request,
	layoutName string,
	templateName string,
	data interface{}) {
	layoutPath := GetTemplatePath(path.Join("layouts", layoutName))
	contentPath := GetTemplatePath(templateName)
	t, err := template.ParseFiles(layoutPath, contentPath)
	if err != nil {
		log.Printf("Unable to find templates. layout:%s, template:%s\n\t%v", layoutPath, contentPath, err)
		http.NotFound(w, r)
		return
	}
	err = t.ExecuteTemplate(w, "layouts/"+layoutName, data)
	if err != nil {
		log.Printf("Failed to execute templates. layout:%s, template:%s\n\t%v", layoutPath, contentPath, err)
		http.NotFound(w, r)
		return
	}
}

// ShowTemplate parses the template and writes the template with the given name to
// the response.
func ShowTemplate(w http.ResponseWriter, r *http.Request, templateName string, data interface{}) {
	tp := GetTemplatePath(templateName)
	t, err := template.ParseFiles(tp)
	if err != nil {
		log.Printf("Unable to find template: %s\n\t%v", tp, err)
		http.NotFound(w, r)
		return
	}
	t.Execute(w, data)
}

var templateDir = determineTemplateDir()
