// Package viewhost is how package ui has a window host a view of the
// platform (ui.NativeView), through a function package native sets as it
// initializes: apps without a native module link no hosting of views.
package viewhost

import "github.com/egoist/mygo/internal/platform"

// Embed has the surface s host view, an NSView* on macOS, as it hosts a
// web view: nil, or returning nil, where it cannot. Main thread only.
var Embed func(s platform.Surface, view uintptr) platform.WebView
