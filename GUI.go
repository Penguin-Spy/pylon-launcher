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

func onImSoCute() {
	fmt.Println("Im sooooooo cute!!")
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

func game_list() []g.Widget {
	items := make([]g.Widget, len(games))

	for i, game := range games {
		items[i] = g.Button(game.name).OnClick(func() {
			fmt.Printf("button for %s clicked: %v\n", game.name, game)
			viewing_game = game
		})
	}

	return items
}

func game_view(game *game) g.Widget {
	return g.Column(
		FullImage(game.id, game.hero),
		g.Label(game.name),
		g.Button("Click Me").OnClick(onClickMe),
		g.Button("I'm so cute").OnClick(onImSoCute),
	)
}

var wnd *g.MasterWindow
var split float32 = 200

type GUIState int

const (
	guistate_loading GUIState = iota
	guistate_game
)

var gui_state GUIState = guistate_loading
var viewing_game *game

func loop() {

	switch gui_state {
	case guistate_loading:
		g.SingleWindow().Layout(
			g.Align(g.AlignCenter).To(g.Label("loading...")),
		)

	case guistate_game:
		g.SingleWindow().Layout(
			g.SplitLayout(g.DirectionVertical, &split,
				g.Column(game_list()...),
				g.Column(game_view(viewing_game)),
			),
		)
	}
}

func GUI_start() {
	wnd = g.NewMasterWindow("Pylon Launcher", 900, 600, 0)

	wnd.Run(loop)
}
