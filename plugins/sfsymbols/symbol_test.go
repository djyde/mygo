package sfsymbols

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// box is an SVG square filling its view box, for a fallback that shows.
var box = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 10 10"><rect width="10" height="10" fill="currentColor"/></svg>`))

// TestSymbolElement checks that a Symbol is named for assistive technology
// and sized as asked, at the theme's font size by default.
func TestSymbolElement(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Row(c).Children(func() {
			Symbol(c, "gearshape")
			Symbol(c, "trash").Size(32).Color(ui.RGB(255, 0, 0)).Weight(Bold).Fill()
		})
	}, 200, 100)
	r, ok := tt.Find("gearshape")
	if !ok || r.W != r.H || r.W <= 0 {
		t.Errorf("gearshape: found %v at %v, want a square", ok, r)
	}
	if r, ok := tt.Find("trash"); !ok || r.W != 32 || r.H != 32 {
		t.Errorf("trash: found %v at %v, want 32 square", ok, r)
	}
}

// TestSymbolFallback checks that a name no system has draws the fallback,
// in the symbol's color.
func TestSymbolFallback(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Fill().Background(ui.RGB(255, 255, 255)).Children(func() {
			Symbol(c, "no.such.symbol").Size(20).Color(ui.RGB(255, 0, 0)).Fallback(box)
		})
	}, 40, 40)
	tt.Frame()
	r, ok := tt.Find("no.such.symbol")
	if !ok {
		t.Fatal("no.such.symbol not found")
	}
	img := tt.Image()
	s := img.Rect.Dx() / 40
	px := img.RGBAAt(int(r.X+r.W/2)*s, int(r.Y+r.H/2)*s)
	if px.R < 200 || px.G > 60 || px.B > 60 {
		t.Errorf("the fallback's middle is %v, want red", px)
	}
}

// TestSymbolBaseline checks that a symbol in a row aligned by baselines
// sits on the text's baseline as SwiftUI's do, the middle of its ink
// within half a DIP of halfway up the capitals, and that a row centering
// its children still centers it on the text's box.
func TestSymbolBaseline(t *testing.T) {
	red, blue := ui.RGB(255, 0, 0), ui.RGB(0, 0, 255)
	for _, align := range []ui.Align{ui.Baseline, ui.Center} {
		tt := ui.NewTester(func(c *ui.Context) {
			ui.Column(c).Fill().Background(ui.RGB(255, 255, 255)).Padding(10, 10, 10, 10).Children(func() {
				ui.Row(c).AlignItems(align).Gap(4).Children(func() {
					Symbol(c, "no.such.symbol").Size(13).Color(blue).Fallback(box)
					ui.Text(c, "HHH").FontSize(13).TextColor(red)
				})
			})
		}, 100, 40)
		tt.SetScale(2)
		tt.Frame()
		sym, _ := tt.Find("no.such.symbol")
		text, _ := tt.Find("HHH")
		if align == ui.Center {
			if d := (sym.Y + sym.H/2) - (text.Y + text.H/2); d > 0.01 || d < -0.01 {
				t.Errorf("in a centered row, the symbol's middle is %v DIPs below the text's", d)
			}
			continue
		}
		// The middle of the ink of a color, in DIPs: the capitals' runs
		// from the cap height down to the baseline.
		img := tt.Image()
		scale := float32(img.Rect.Dx()) / 100
		middle := func(r ui.Rect, red bool) float32 {
			top, bottom := -1, -1
			for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
				for x := int(r.X * scale); x < int((r.X+r.W)*scale); x++ {
					px := img.RGBAAt(x, y)
					if (red && px.R > 128 && px.B < 128) || (!red && px.B > 128 && px.R < 128) {
						if top < 0 {
							top = y
						}
						bottom = y + 1
						break
					}
				}
			}
			if top < 0 {
				return -1
			}
			return float32(top+bottom) / 2 / scale
		}
		capMid, symMid := middle(text, true), middle(sym, false)
		if capMid < 0 {
			t.Skip("the text engine draws no glyphs")
		}
		if d := symMid - capMid; d > 0.5 || d < -0.5 {
			t.Errorf("in a baseline row, the symbol's middle is %v DIPs below the capitals', want within 0.5", d)
		}
	}
}
