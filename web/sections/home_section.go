package websections

import "github.com/maxence-charriere/go-app/v10/pkg/app"

// HomeSection is a component that renders the home section.
type HomeSection struct {
	app.Compo
}

// Render renders the home section.
func (s *HomeSection) Render() app.UI {
	return app.Div().Class("container").Body()
}
