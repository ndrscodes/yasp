package templates

import (
	"embed"
	"errors"
	"html/template"
	"time"
)

//go:embed layout pages
var templates embed.FS

type TemplateRegistry struct {
	templates map[string]*template.Template
	base      *template.Template
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

func (tr *TemplateRegistry) Register(page string) (*template.Template, error) {
	if tr.base == nil {
		return nil, errors.New("base template has not been registered")
	}

	b, err := tr.base.Clone()
	if err != nil {
		return nil, err
	}

	t, err := b.ParseFS(templates, "pages/"+page+"/*")
	if err != nil {
		return nil, err
	}

	tr.templates[page] = t
	return t, nil
}

func NewTemplateRegistry() (t *TemplateRegistry, err error) {
	b, err := template.ParseFS(templates, "layout/*")
	if err != nil {
		return nil, err
	}

	return &TemplateRegistry{
		templates: make(map[string]*template.Template),
		base:      b,
		updated:   time.Now(),
	}, nil
}
