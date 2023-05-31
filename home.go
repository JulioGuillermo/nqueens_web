package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"

	"github.com/julioguillermo/nqueens_web/nqueens"
)

type Home struct {
	app.Compo

	drawer      app.Value
	algDropdown app.Value

	running bool

	N       int
	Alg     string
	PopSize int
	Surv    int
	MutRate float64
	MC      bool

	Threads int

	Generation  uint64
	Progress    float64
	Duration    time.Duration
	Errors      int
	TotalErrors int

	exporting bool
	output    string
}

func (p *Home) OnMount(ctx app.Context) {
	p.N = 8
	p.Alg = "Genetic"
	p.PopSize = 1000
	p.Surv = 10
	p.MutRate = 0.2
	p.MC = true
	p.Threads = 10

	ctx.Defer(p.initJSComponents)
}

func (p *Home) initJSComponents(ctx app.Context) {
	dr := app.Window().GetElementByID("drawer")
	p.drawer = app.Window().Get("M").Get("Sidenav").Call("init", dr)

	rbtn := app.Window().GetElementByID("algBtn")
	p.algDropdown = app.Window().Get("M").Get("Dropdown").Call("init", rbtn)
}

func (p *Home) OnDismount() {
	if p.drawer != nil {
		p.drawer.Call("destroy")
	}
	if p.algDropdown != nil {
		p.algDropdown.Call("destroy")
	}
}

func (p *Home) Render() app.UI {
	return app.Div().Body(
		app.Div().Class("navbar-fixed").Body(
			app.Nav().Class("primary top-nav").Body(
				app.Div().
					Class("nav-wrapper").
					Style("padding-left", "1cm").
					Style("padding-right", "1cm").
					Body(
						app.A().
							Href("#").
							Class("sidenav-trigger").
							DataSet("target", "drawer").
							Body(
								app.I().Class("material-icons").Body(app.Text("menu")),
							),
						app.A().
							Href("/").
							Class("brand-logo").
							Text("NQueens"),
						app.If(
							!p.running,
							app.Ul().Class("right").Body(
								app.Li().Body(
									app.A().Href("#").OnClick(func(ctx app.Context, e app.Event) {
										p.RunGANQ(ctx)
									}).Body(
										app.I().Class("material-icons left").Text("play_arrow"),
										app.Text("Start"),
									),
								),
							),
						).Else(
							app.Ul().Class("right").Body(
								app.Li().Body(
									app.A().Href("#").OnClick(func(ctx app.Context, e app.Event) {
										p.running = false
									}).Body(
										app.I().Class("material-icons left").Text("stop"),
										app.Text("Stop"),
									),
								),
							),
						),
					),
			),
		),
		app.Ul().ID("drawer").Class("sidenav sidenav-fixed").Body(
			app.Li().Body(
				app.Div().
					Class("user-view").
					Style("position", "relative").
					Style("height", "200px").
					Body(
						app.Div().Class("background").
							Style("height", "100%").
							Body(
								app.Img().
									Class("responsive-img").
									Src("/web/assets/background.jpeg"),
							),
						app.Div().
							Style("position", "absolute").
							Style("left", "2cm").
							Style("bottom", "0.3cm").
							Class("container").
							Body(
								app.Img().
									Class("responsive-img").
									Style("width", "1cm").
									Src("/web/assets/icon.png"),
							),
						app.Div().
							Style("position", "absolute").
							Style("right", "0").
							Style("top", "0").
							Class("container").
							Body(
								app.Span().ID("title").Text("N Queens"),
							),
					),
			),

			app.Li().Class("divider"),

			app.Li().Body(
				app.Div().Class("input-field").Body(
					app.Input().
						ID("n").
						Type("number").
						Min("4").
						Value(p.N).
						OnChange(p.ValueTo(&p.N)),
					app.Label().Class("active").For("n").Text("N Queens"),
				),
			),

			/*app.Li().Class("divider"),

			app.Li().Body(
				app.A().
					ID("algBtn").
					Href("#").
					Class("dropdown-trigger w100").
					DataSet("target", "alg_dropdown").
					Body(
						app.I().
							Class("material-icons right").
							Text("arrow_drop_down"),
						app.Text(p.Alg),
					),
				app.Ul().ID("alg_dropdown").Class("dropdown-content").Body(
					app.Li().Body(
						app.A().OnClick(func(ctx app.Context, e app.Event) {
							p.Alg = "Genetic"
						}).Href("#").Text("Genetic"),
					),
					app.Li().Body(
						app.A().OnClick(func(ctx app.Context, e app.Event) {
							p.Alg = "Random"
						}).Href("#").Text("Random"),
					),
					app.Li().Body(
						app.A().OnClick(func(ctx app.Context, e app.Event) {
							p.Alg = "Simple"
						}).Href("#").Text("Simple"),
					),
					app.Li().Body(
						app.A().OnClick(func(ctx app.Context, e app.Event) {
							p.Alg = "Heuristic"
						}).Href("#").Text("Heuristic"),
					),
					app.Li().Body(
						app.A().OnClick(func(ctx app.Context, e app.Event) {
							p.Alg = "Fast Heuristic"
						}).Href("#").Text("Fast Heuristic"),
					),
				),
			),*/

			app.If(
				p.Alg == "Genetic",
				app.Li().Class("divider"),

				app.Li().Body(
					app.Div().Class("input-field").Body(
						app.Input().
							ID("Threads").
							Type("number").
							Min("1").
							Value(p.Threads).
							OnChange(p.ValueTo(&p.Threads)),
						app.Label().Class("active").For("Threads").Text("Threads"),
					),
				),
				app.Li().Body(
					app.Div().Class("input-field").Body(
						app.Input().
							ID("PopSize").
							Type("number").
							Min("50").
							Value(p.PopSize).
							OnChange(p.ValueTo(&p.PopSize)),
						app.Label().Class("active").For("PopSize").Text("Population size"),
					),
				),
				app.Li().Body(
					app.Div().Class("input-field").Body(
						app.Input().
							ID("Surv").
							Type("number").
							Min("5").
							Max(p.PopSize/2).
							Value(p.Surv).
							OnChange(p.ValueTo(&p.Surv)),
						app.Label().Class("active").For("Surv").Text("Survivors"),
					),
				),

				app.Li().Body(
					app.Div().Class("input-field").Body(
						app.Input().
							ID("MutRate").
							Type("number").
							Min(0).
							Max(1).
							Value(p.MutRate).
							OnChange(p.ValueTo(&p.MutRate)),
						app.Label().Class("active").For("MutRate").Text("Mutation rate"),
					),
				),

				app.Li().Body(
					app.Div().Class("switch mb4").Body(
						app.Label().Body(
							app.Input().
								Type("checkbox").
								Checked(p.MC).
								OnChange(func(ctx app.Context, e app.Event) {
									p.MC = e.Value.Get("target").
										Get("checked").
										Bool()
								}),
							app.Span().Class("lever"),
							app.Text("Mapped crossover"),
						),
					),
				),
			),
		),

		app.Main().Body(
			app.If(
				p.running,
				app.Div().Class("container").Body(
					app.Div().Class("row").Body(
						app.Div().Class("col s12").Body(
							app.Div().Class("card").Body(
								app.Div().Class("card-content").Body(
									app.Div().Class("container").Body(
										app.If(
											p.exporting,
											app.Div().
												Class("card-title primary-text").
												Text("Exporting SVG"),
											app.Div().Class("progress secondary").Body(
												app.Div().
													Class("indeterminate primary"),
											),
										).Else(
											app.Div().
												Class("card-title primary-text").
												Text("Running"),
											app.Div().Class("progress secondary").Body(
												app.Div().
													Class("determinate primary").
													Style("width", fmt.Sprintf("%0.2f%%", p.Progress)),
											),
										),
										app.Table().Body(
											app.Tr().Body(
												app.Td().Text("Generation"),
												app.Td().Text(p.Generation),
											),
											app.Tr().Body(
												app.Td().Text("Duration"),
												app.Td().
													Text(fmt.Sprintf("%0.4f", p.Duration.Seconds())),
											),
											app.Tr().Body(
												app.Td().Text("Errors"),
												app.Td().
													Text(fmt.Sprintf("%d / %d", p.Errors, p.TotalErrors)),
											),
										),
									),
								),
							),
						),
					),
				),
			).Else(
				app.Div().Class("container").Body(
					app.Div().Class("row").Body(
						app.Div().Class("col s12").Body(
							app.Raw(p.output),
						),
					),
				),
			),
		),
	)
}

func (p *Home) RunGANQ(ctx app.Context) {
	p.Progress = 0
	p.running = true
	p.exporting = false
	ga := nqueens.NewGA(p.N, p.PopSize, p.Surv, p.MutRate, p.MC, p.Threads)

	var m sync.Mutex
	ctx.Async(func() {
		for p.running && ga.GetBestError() > 0 {
			m.Lock()
			ga.NextGen()
			ctx.Dispatch(func(ctx app.Context) {
				p.Generation, p.Progress, p.Duration, p.Errors, p.TotalErrors = ga.Info()
				app.Logf(
					"Gen: %d\nPro: %f\nTime: %v\nErr: %d / %d",
					p.Generation,
					p.Progress,
					p.Duration,
					p.Errors,
					p.TotalErrors,
				)
				time.Sleep(time.Millisecond)
				m.Unlock()
			})
		}

		ctx.Dispatch(func(ctx app.Context) {
			p.exporting = true
		})
		output, _ := ga.Save()
		ctx.Dispatch(func(ctx app.Context) {
			p.output = output
			p.running = false
		})
	})
}
