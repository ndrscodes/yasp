package templates

import (
	"html/template"
	"time"
)

type TemplateRegistry struct {
	template *template.Template
	path     string
	updated  time.Time
}

// retrieve the contained template.
// this struct merely serves as a single point of truth for templates. It may later be rewritten to act as a cache
// with a timeout.
func (tr *TemplateRegistry) Get() (*template.Template, error) {
	return tr.template, nil
}

// reloads the contained template
func (tr *TemplateRegistry) Reload() (t *template.Template, err error) {
	t, err = template.ParseGlob(tr.path)
	if err != nil {
		return nil, err
	}

	tr.template = t
	tr.updated = time.Now()

	return t, nil
}

func NewTemplateRegistry(path string) (t *TemplateRegistry, err error) {
	tmp, err := template.ParseGlob(path)
	if err != nil {
		return nil, err
	}

	return &TemplateRegistry{
		template: tmp,
		path:     path,
		updated:  time.Now(),
	}, nil
}
