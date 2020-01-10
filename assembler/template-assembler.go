package assembler

import (
	"text/template"

	rice "github.com/GeertJohan/go.rice"
)

// AssembleTemplate ...
func AssembleTemplate(box *rice.Box, base string, view string) (*template.Template, error) {
	base, err := box.String(base)
	if err != nil {
		return nil, err
	}

	content, err := box.String(view)
	if err != nil {
		//log.Panic(err.Error())
		return nil, err
	}

	x, err := template.New("base").Parse(base)
	if err != nil {
		//log.Panic(err.Error())
		return nil, err
	}

	x.New("content").Parse(content)
	if err != nil {
		//log.Panic(err.Error())
		return nil, err
	}
	return x, nil
}
