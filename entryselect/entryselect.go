package entryselect

import (
	"fmt"
	"reflect"
	"runtime"

	survey "github.com/AlecAivazis/survey/v2"
)

type entry struct {
	Name string
	Func func()
}

var mainEntries map[string]*entry

func Add(fns ...func()) {
	for _, fn := range fns {
		e := &entry{
			Name: getFunctionName(fn),
			Func: fn,
		}
		if _, ok := mainEntries[e.Name]; ok {
			panic("duplicated entry name:" + e.Name)
		}
		if mainEntries == nil {
			mainEntries = make(map[string]*entry)
		}
		mainEntries[e.Name] = e
	}
}

func Execute() {
	entryOptions := getEntryOptions()
	if len(entryOptions) <= 0 {
		panic("zero entries, please add one")
	}
	qs := []*survey.Question{
		{
			Name: "entry",
			Prompt: &survey.Select{
				Message: "Choose a main function:",
				Options: entryOptions,
				VimMode: true,
			},
		},
	}
	ans := struct {
		Entry string
	}{}
	if err := survey.Ask(qs, &ans); err != nil {
		fmt.Println(err)
		return
	}
	fn, ok := mainEntries[ans.Entry]
	if !ok {
		panic("entry not found:" + ans.Entry)
	}
	fn.Func()
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func getEntryOptions() []string {
	var entryOptions []string
	for _, e := range mainEntries {
		entryOptions = append(entryOptions, e.Name)
	}
	return entryOptions
}
