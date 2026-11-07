// Native-module is an app of native UI that calls native modules, Go
// interfaces implemented in Swift: package language in modules/language,
// with Apple's NaturalLanguage framework, tells the language and the
// sentences of the text typed; package places in modules/places shows a
// map of MapKit, a view of the platform, in the window, and the system's
// share menu.
//
//	go run ./cmd/mygo dev examples/native-module
//
// mygo dev and mygo build generate the modules and build their Swift
// packages into the app (see docs/native-modules.md).
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/examples/native-module/modules/language"
	"github.com/egoist/mygo/examples/native-module/modules/places"
	"github.com/egoist/mygo/ui"
)

type app struct {
	win      *mygo.Window
	analyzer language.Analyzer
	err      error
	text     string

	// The map, a view of the platform, made once and shown by the view.
	m       places.Map
	sharing places.Sharing
	center  places.Place
	moved   io.Closer
}

func (a *app) view(c *ui.Context) {
	t := c.Theme()
	ui.Row(c).Fill().Padding(24).Gap(24).Children(func() {
		ui.Column(c).Grow(1).AlignSelf(ui.Stretch).Gap(12).Children(func() { a.text_(c, t) })
		if a.m == nil {
			return
		}
		ui.Column(c).Grow(1).AlignSelf(ui.Stretch).Gap(12).Children(func() {
			// What the view paints after the map shows over it.
			ui.Box(c).Grow(1).Radius(12).Children(func() {
				ui.NativeView(c, a.m).Fill().Radius(12)
				ui.Textf(c, "%.3f, %.3f", a.center.Latitude, a.center.Longitude).
					Absolute().Right(8).Bottom(8).Padding(8).Background(t.Surface).Radius(6)
			})
			ui.Row(c).Gap(8).Children(func() {
				share := ui.Button(c.Key("share"), "Share")
				if share.Clicked() {
					a.sharing.Share(a.win, share.Bounds(), a.text)
				}
				if ui.Button(c.Key("paris"), "Paris").Clicked() {
					a.m.Show(places.Place{Latitude: 48.8566, Longitude: 2.3522}, 0.2)
				}
				if ui.Button(c.Key("tokyo"), "Tokyo").Clicked() {
					a.m.Show(places.Place{Latitude: 35.6762, Longitude: 139.6503}, 0.3)
				}
			})
		})
	})
}

func (a *app) text_(c *ui.Context, t *ui.Theme) {
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
}

func main() {
	a := &app{text: "Bonjour tout le monde. Comment allez-vous ?"}
	// errors.ErrUnsupported where a module has no implementation.
	a.analyzer, a.err = language.NewAnalyzer()
	mygo.App.WhenReady(func() {
		// Views are made on the main thread, once the app runs.
		if m, err := places.NewMap(); err == nil {
			a.m = m
			a.sharing, _ = places.NewSharing()
			m.Show(places.Place{Latitude: 48.8566, Longitude: 2.3522}, 0.2)
			a.center = m.Center()
			a.moved = m.OnMove(func(p places.Place) {
				a.center = p
				a.win.Invalidate()
			})
		}
		a.win = mygo.NewWindow(mygo.WindowOptions{
			Title:   "Native modules",
			Width:   900,
			Height:  520,
			Content: ui.View(a.view),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
