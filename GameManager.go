/* GameManager.go © Penguin_Spy 2024
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

// https://stackoverflow.com/a/49595208
func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

type game struct {
	id       string      // the internal ID of this game
	name     string      // the user-facing title of this game
	hero     image.Image // the "hero" image, shown in the background of the game's page
	pluginId string      // the plugin that defined this game
}

func Game(id string, name string, heroFile string, pluginId string) *game {
	image, err := getImageFromFilePath(heroFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[GameManager] loading image %q failed\n", heroFile)
		panic(err)
	}
	return &game{
		id,
		name,
		image,
		pluginId,
	}
}

var games []*game

func RegisterGame(game *game) {
	games = append(games, game)
}
