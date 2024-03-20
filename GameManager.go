/* GameManager.go © Penguin_Spy 2024
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
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
	id     string      // the internal ID of this game
	name   string      // the user-facing title of this game
	hero   image.Image // the "hero" image, shown in the background of the game's page
	plugin *plugin     // the plugin that defined this game
}

var games []*game

// creates a game object that represents a game. does not add it to the registered games list.
func DefineGame(id string, name string, heroFile string, plugin *plugin) (*game, error) {
	image, err := getImageFromFilePath(filepath.Join(getPluginDir(plugin.id), heroFile))
	if err != nil {
		log.Printf("[GameManager] loading image %q failed\n", heroFile)
		return nil, err
	}
	game := &game{
		id,
		name,
		image,
		plugin,
	}
	return game, nil
}

// registers a game that was previously defined with DefineGame
func RegisterGame(game *game) {
	games = append(games, game)
}
