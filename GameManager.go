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
	id   string
	name string
	hero image.Image
}

func Game(id string, name string, heroFile string) *game {
	image, err := getImageFromFilePath(heroFile)
	if err != nil {
		fmt.Printf("loading image %q failed\n", heroFile)
		panic(err)
	}
	return &game{
		id,
		name,
		image,
	}
}

var games []*game

func RegisterGame(game *game) {
	games = append(games, game)
}