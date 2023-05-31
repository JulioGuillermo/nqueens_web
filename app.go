package main

import (
	"os"
	"path"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func GetApp() *app.Handler {
	app.Route("/", &Home{})
	app.RouteWithRegexp("/.*", &Home{})
	app.RunWhenOnBrowser()

	a := &app.Handler{
		Icon: app.Icon{
			Default:    "web/assets/icon_192.png",
			Large:      "web/assets/icon.png",
			AppleTouch: "web/assets/icon_192.png",
		},
		Title:           "NQueens",
		Name:            "NQueens",
		ShortName:       "NQueens",
		BackgroundColor: "#FFFFFF",
		ThemeColor:      "#0088ff",
		Version:         "0.0.1",
		Image:           "web/assets/icon.png",
		Author:          "Julio Guillermo Mayo Vidal",
		Description:     "NQueens web",
		Styles: []string{
			"web/assets/material-icon.css",
			"web/assets/materialize.min.css",
			"web/assets/style.css",
		},
		Scripts: []string{
			"web/assets/materialize.min.js",
			"web/assets/particles.min.js",
			"web/assets/particles_init.js",
		},
		CacheableResources: GetResources("web/"),
	}

	return a
}

func GetResources(p string) []string {
	queue := []string{p}
	res := []string{}

	for len(queue) > 0 {
		p = queue[0]
		queue = queue[1:]

		elements, err := os.ReadDir(p)
		if err != nil {
			continue
		}

		for _, e := range elements {
			if e.IsDir() {
				queue = append(queue, path.Join(p, e.Name()))
			} else {
				res = append(res, path.Join(p, e.Name()))
			}
		}
	}

	return res
}
