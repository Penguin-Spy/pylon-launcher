/* main.go © Penguin_Spy 2024
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	g "github.com/AllenDang/giu"
	lua "github.com/yuin/gopher-lua"
)

func main() {
	go func() {
		L := lua.NewState()
		defer L.Close()
		launcher_open_lua_lib(L)
		if err := L.DoFile("testplugin.lua"); err != nil {
			panic(err)
		}
		// once all plugins have loaded
		viewing_game = games[0]
		gui_state = guistate_game
		g.Update()
	}()

	GUI_start()
}
