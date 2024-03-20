/* lua_launcher.go © Penguin_Spy 2024
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	lua "github.com/yuin/gopher-lua"
)

func OpenLauncher(L *lua.LState) int {
	// register functions to the table
	mod := L.RegisterModule("launcher", map[string]lua.LGFunction{
		"define_game": define_game,
		"on_play":     on_play,
	})
	// register other stuff
	L.SetField(mod, "name", lua.LString("value"))
	L.Push(mod)
	return 1
}

func define_game(L *lua.LState) int {
	data := L.CheckTable(1)

	var id, name, heroPath lua.LString
	var ok bool

	if id, ok = L.GetField(data, "id").(lua.LString); !ok {
		L.ArgError(1, "invalid game data table: id")
		return 0
	}
	if name, ok = L.GetField(data, "name").(lua.LString); !ok {
		L.ArgError(1, "invalid game data table: name")
		return 0
	}
	if heroPath, ok = L.GetField(data, "hero").(lua.LString); !ok {
		L.ArgError(1, "invalid game data table: hero")
		return 0
	}

	pluginId := L.GetField(L.Get(lua.RegistryIndex), "plugin_id").String()
	plugin := plugins[pluginId]
	if game, err := DefineGame(id.String(), name.String(), heroPath.String(), plugin); err != nil {
		L.RaiseError("failed to define game - %v", err)
	} else {
		plugin.games = append(plugin.games, game)
	}

	return 0
}

func on_play(L *lua.LState) int {
	callback := L.CheckFunction(1)
	L.SetField(L.Get(lua.RegistryIndex), "on_play", callback)
	return 0
}
