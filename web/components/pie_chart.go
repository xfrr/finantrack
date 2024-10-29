package wasmcomponents

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type ChartData struct {
	Labels   []string
	Datasets []ChartDataset
}

type ChartDataset struct {
	Label           string
	Data            []float64
	BackgroundColor []string
}

// PieChart is a component that uses chart.js to render a pie chart.
type PieChart struct {
	app.Compo

	Data any
}

// Render returns the component's HTML markup.
func (p *PieChart) Render() app.UI {
	return app.Div().
		Class("row").
		Body(
			app.Div().
				Class("col-12").
				Body(
					app.Canvas().
						ID("pie-chart").
						Class("chartjs-render-monitor").
						Style("height", "400px"),
				),
		)
}

// JS is called after the component is rendered.
func (p *PieChart) JS() {
	app.Window().Call("renderPieChart", p.Data)
}

// JS is a JavaScript function that uses chart.js to render a pie chart.
const JS = `
function renderPieChart(data) {
	var ctx = document.getElementById('pie-chart').getContext('2d');
	var myPieChart = new Chart(ctx, {
		type: 'pie',
		data: data,
		options: {
			responsive: true,
			maintainAspectRatio: false,
		},
	});
}
`

// Styles returns the component's CSS.
func (p *PieChart) Styles() string {
	return `
		.chartjs-render-monitor {
			width: 100%;
		}
	`
}

// Styles is the component's CSS.
const Styles = `
	.chartjs-render-monitor {
		width: 100%;	
	}
`
