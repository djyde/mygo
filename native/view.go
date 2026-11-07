package native

import (
	"github.com/egoist/mygo/internal/platform"
	"github.com/egoist/mygo/internal/viewhost"
)

// The windows of apps that use native modules host their views, which
// ui.NativeView shows.
func init() { viewhost.Embed = embed }

func embed(s platform.Surface, view uintptr) platform.WebView {
	vs, ok := s.(platform.ViewSurface)
	if !ok {
		return nil
	}
	e, err := vs.EmbedView(view)
	if err != nil {
		return nil
	}
	return e
}
