package websections

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	wasmcomponents "github.com/xfrr/finantrack/web/components"
)

const Path = "/assets/[a-zA-Z0-9].*"

type Asset struct {
	Name        string
	Description string
	Type        string
	Money       float64
}

// AssetsSection is a component that renders the assets section.
type AssetsSection struct {
	app.Compo

	Assets []*Asset
}

func (a *AssetsSection) renderFloatingActionBar() app.UI {
	return app.Div().
		Style("position", "fixed").
		Style("bottom", "1.5rem").
		Style("right", "3rem").
		Style("z-index", "1000").
		Style("display", "flex").
		Body(
			app.Button().
				Class("btn btn-primary").
				Style("border-radius", "50%").
				Style("display", "flex").
				Style("align-items", "center").
				Style("justify-content", "center").
				Style("width", "4rem").
				Style("height", "4rem").
				Body(app.I().Class("fas fa-plus")),
		)
}

func (a *AssetsSection) renderAssetList() app.UI {
	return app.Div().
		Class("container").
		Body(
			app.If(a.Assets != nil, func() app.UI {
				return app.Range(a.Assets).Slice(func(i int) app.UI {
					asset := a.Assets[i]
					return app.Div().Class("card mb-3").Body(
						app.Div().Class("card-body").Body(
							&wasmcomponents.ScrollingListItem{
								Title:    asset.Name,
								Subtitle: asset.Description,
								Actions: app.Span().Body(
									app.Text("Money: $"),
									app.Text(app.FormatString("%.2f", asset.Money)),
								),
							},
						),
					)
				})
			}),
		)
}

func (a *AssetsSection) renderAssetChart() app.UI {
	return &wasmcomponents.PieChart{
		Data: wasmcomponents.ChartData{
			Labels: []string{"Red", "Blue", "Yellow"},
			Datasets: []wasmcomponents.ChartDataset{
				{
					Data:            []float64{300, 50, 100},
					BackgroundColor: []string{"#FF6384", "#36A2EB", "#FFCE56"},
				},
			},
		},
	}
}

// Render renders the assets section.
func (a *AssetsSection) Render() app.UI {
	return app.Div().
		Class("container-fluid h-100 w-100").
		Body(
			app.Div().
				Class("row").
				Body(
					app.Div().
						Class("col-12").
						Body(
							app.H1().Text("Assets"),
							a.renderAssetChart(),
						),
				),
			a.renderAssetList(),
			a.renderFloatingActionBar(),
		)
}
