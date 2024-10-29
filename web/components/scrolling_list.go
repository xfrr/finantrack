package wasmcomponents

import "github.com/maxence-charriere/go-app/v10/pkg/app"

// ScrollingList is a wasm component that renders a list of items.
// Items can be scrolled horizontally.
// This component uses Bootstrap.
type ScrollingList struct {
	app.Compo

	Items []ScrollingListItem
}

func (l *ScrollingList) OnMount(ctx app.Context) {
	l.Items = []ScrollingListItem{
		{
			Title:    "Item 1",
			Subtitle: "This is the first item",
		},
		{
			Title:    "Item 2",
			Subtitle: "This is the second item",
		},
		{
			Title:    "Item 3",
			Subtitle: "This is the third item",
		},
	}
}

// Render renders the list using the items and Bootstrap.
func (l ScrollingList) Render() app.UI {
	return app.Div().Class("d-flex overflow-auto").Body(
		app.Range(l.Items).Slice(func(i int) app.UI {
			return l.Items[i].Render()
		}),
	)
}
