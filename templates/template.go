package template

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/template"

	changecase "github.com/ku/go-change-case"
	"github.com/sirupsen/logrus"
)

//go:generate go-bindata -pkg=template -ignore=.go -nomemcopy  tmpl/... fonts/...

var engine Engine

type Engine interface {
	init()
	Execute(name string, model interface{}) (string, error)
	ExecuteString(data string, model interface{}) (string, error)
	MustAssetString(name string) string
}

type DefaultEngine struct {
	t *template.Template
}

func funcMap() template.FuncMap {
	return template.FuncMap{
		"last": func(x int, a interface{}) bool {
			return x == reflect.ValueOf(a).Len()-1
		},
		"toParam": func(s string) string {
			return changecase.Param(s)
		},
		"removeHash": func(s string) string {
			return strings.Replace(s, "#", "", -1)
		},
		"divide": func(a, b int) float64 {
			return float64(a) / float64(b)
		},
		"times": func(a int, b float64) float64 {
			return float64(a) * b
		},
	}
}
func NewEngine() Engine {
	if engine == nil {
		engine = &DefaultEngine{}
		engine.init()
	}
	return engine
}
func (e *DefaultEngine) init() {
	e.t = template.New("default")
	e.t.Funcs(funcMap())
	for n, v := range _bintree.Children["tmpl"].Children["partials"].Children {
		a, _ := v.Func()
		_, err := e.t.Parse(
			fmt.Sprintf(
				"{{define \"%s\"}} %s {{end}}",
				strings.Replace(n, ".tmpl", "", 1),
				string(a.bytes),
			),
		)
		if err != nil {
			logrus.Panic(err)
		}
	}
}

func (e *DefaultEngine) Execute(name string, model interface{}) (string, error) {
	d, err := Asset(fmt.Sprintf("tmpl/%s.tmpl", name))
	if err != nil {
		logrus.Panic(err)
	}
	tmp, err := e.t.Parse(string(d))
	if err != nil {
		logrus.Panic(err)
	}
	ret := bytes.NewBufferString("")
	err = tmp.Execute(ret, model)
	return ret.String(), err
}

func (e *DefaultEngine) ExecuteString(data string, model interface{}) (string, error) {
	tmp, err := e.t.Parse(data)
	if err != nil {
		logrus.Panic(err)
	}
	ret := bytes.NewBufferString("")
	err = tmp.Execute(ret, model)
	return ret.String(), err
}

func (e *DefaultEngine) MustAssetString(name string) string {
	return MustAssetString(name)
}
