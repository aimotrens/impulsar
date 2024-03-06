package engine

import (
	"strings"
	"text/template"

	"github.com/aimotrens/impulsar/model"
)

func (e *Engine) ExpandVarsWithTemplateEngine(script string, j *model.Job) string {
	vars := e.aggregateEnvVars(j)

	t := template.Must(template.New("template").
		Funcs(template.FuncMap{
			"iterate": iterate,
			"split":   strings.Split,
		}).
		Option("missingkey=zero").
		Parse(script),
	)

	buf := new(strings.Builder)

	err := t.Execute(buf, vars)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

func iterate(from, to int) []int {
	var result []int
	for i := from; i <= to; i++ {
		result = append(result, i)
	}
	return result
}
