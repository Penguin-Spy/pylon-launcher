/* GUI.go © Penguin_Spy 2024
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"fmt"
	"image"

	g "github.com/AllenDang/giu"
)

func onClickMe() {
	fmt.Println("Hello world!")
}

type FullImageWidget struct {
	id    string
	image image.Image
}

func FullImage(id string, image image.Image) *FullImageWidget {
	return &FullImageWidget{
		id:    id,
		image: image,
	}
}

func (w *FullImageWidget) Build() {
	width, _ := g.GetAvailableRegion()
	bounds := w.image.Bounds()
	height := float32(bounds.Dy()) * (width / float32(bounds.Dx()))

	g.ImageWithRgba(w.image).
		ID("FullImage##"+w.id).
		Size(width, height).
		Build()
}

func gameList() []g.Widget {
	items := make([]g.Widget, len(games))

	for i, game := range games {
		items[i] = g.Button(game.name).OnClick(func() {
			fmt.Printf("button for %s clicked: %v\n", game.name, game)
			viewingGame = game
		})
	}

	return items
}

func gameView(game *game) g.Widget {
	return g.Column(
		FullImage(game.id, game.hero),
		g.Label(game.name).Font(fonts["header"]),
		g.Button("Click Me").OnClick(onClickMe),
		g.Button("I'm so cute").OnClick(func() {
			callOnPlay(game)
		}),
	)
}

var wnd *g.MasterWindow
var split float32 = 200

type GUIState int

const (
	guiStateLoading GUIState = iota
	guiStateGame
)

var guiState GUIState = guiStateLoading
var viewingGame *game

func loop() {
	switch guiState {
	case guiStateLoading:
		g.SingleWindow().Layout(
			g.Align(g.AlignCenter).To(g.Label("loading...").Font(fonts["header"])),
		)

	case guiStateGame:
		g.SingleWindow().Layout(
			g.SplitLayout(g.DirectionVertical, &split,
				g.Column(gameList()...),
				g.Column(gameView(viewingGame)),
			),
		)
	}
}

var fonts = make(map[string]*g.FontInfo)

func GUI_start() {
	wnd = g.NewMasterWindow("Pylon Launcher", 900, 600, 0)
	g.Context.FontAtlas.SetDefaultFont("fonts/Montserrat-Regular.ttf", 18)
	fonts["header"] = g.Context.FontAtlas.AddFont("fonts/MontserratAlternates-Regular.ttf", 48)
	wnd.Run(loop)
}
