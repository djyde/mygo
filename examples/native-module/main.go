// Native-module is an app of native UI that calls a native module: the
// package language in modules/language, whose interface Go declares and
// Swift implements with Apple's NaturalLanguage framework. Type text, and
// the app tells its language and its sentences.
//
//	go run ./cmd/mygo dev examples/native-module
//
// mygo dev and mygo build generate the module and build its Swift package
// into the app (see docs/native-modules.md).
package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/examples/native-module/modules/language"
	"github.com/egoist/mygo/ui"
)

type app struct {
	analyzer language.Analyzer
	err      error
	text     string
}

func (a *app) view(c *ui.Context) {
	t := c.Theme()
	ui.Column(c).Fill().Padding(24).Gap(12).Children(func() {
		ui.TextArea(c.Key("text"), &a.text).Placeholder("Type in any language…").Height(120).AutoFocus()
		if a.err != nil {
			ui.Text(c, a.err.Error()).TextColor(t.TextMuted)
			return
		}
		lang, err := a.analyzer.DetectLanguage(a.text)
		switch {
		case errors.Is(err, language.ErrUndetermined):
			ui.Text(c, "Language: undetermined").TextColor(t.TextMuted)
		case err != nil:
			ui.Text(c, err.Error())
		default:
			ui.Textf(c, "Language: %s", lang).Bold()
		}
		n := 0
		for tok, err := range a.analyzer.Tokens(context.Background(), a.text, language.Sentence) {
			if err != nil {
				break
			}
			n++
			ui.Text(c.Key(fmt.Sprint("s", n)), fmt.Sprintf("%d. %s", n, tok.Text))
		}
	})
}

func main() {
	a := &app{text: "Bonjour tout le monde. Comment allez-vous ?"}
	// errors.ErrUnsupported where the module has no implementation.
	a.analyzer, a.err = language.NewAnalyzer()
	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{
			Title:   "Native module",
			Width:   520,
			Height:  420,
			Content: ui.View(a.view),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
