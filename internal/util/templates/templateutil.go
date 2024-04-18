package templates

import (
	"errors"
	"html/template"
	"io"
	"os"
	"time"
)

type TemplateRegistry struct {
	templates map[string]*template.Template
	base      *template.Template
	path      string
	updated   time.Time
}

// retrieve the contained template.
// this struct merely serves as a single point of truth for templates. It may later be rewritten to act as a cache
// with a timeout.
func (tr *TemplateRegistry) Get(name string) (*template.Template, error) {
	t := tr.templates[name]

	if t == nil {
		return nil, errors.New("template not found: " + name)
	}

	return t, nil
}

// reloads the contained template
func (tr *TemplateRegistry) Reload(name string) (*template.Template, error) {
	t := tr.templates[name]
	if t == nil {
		return nil, errors.New("template not found: " + name)
	}

	delete(tr.templates, name)
	t, err := tr.Register(name)

	return t, err
}

func (tr *TemplateRegistry) Register(page string) (*template.Template, error) {
	if tr.base == nil {
		return nil, errors.New("base template has not been registered.")
	}

	b, err := tr.base.Clone()
	if err != nil {
		return nil, err
	}

	t, err := b.ParseGlob(tr.path + "/pages/" + page + "/*")
	if err != nil {
		return nil, err
	}

	tr.templates[page] = t
	return t, nil
}

func NewTemplateRegistry(path string) (t *TemplateRegistry, err error) {
	b, err := template.ParseGlob(path + "/layout/*")
	if err != nil {
		return nil, err
	}

	com := path + "/common"
	f, err := os.Open(com)
	if err != nil {
		return nil, err
	}

	_, err = f.Readdirnames(1)
	if err != io.EOF {
		return nil, err
	}
	if err == nil {
		b, err = b.ParseGlob(com + "/*")
		if err != nil {
			return nil, err
		}
	}

	return &TemplateRegistry{
		templates: make(map[string]*template.Template),
		base:      b,
		path:      path,
		updated:   time.Now(),
	}, nil
}
