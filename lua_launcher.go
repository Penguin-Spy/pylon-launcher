/* lua_launcher.go © Penguin_Spy 2024
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"fmt"
	"time"

	lua "github.com/yuin/gopher-lua"
)

func sleep(L *lua.LState) int {
	duration := L.CheckNumber(1)
	time.Sleep(time.Second * time.Duration(duration))
	return 0
}

func launcher_open_lua_lib(L *lua.LState) int {
	// register functions to the table
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"define_game": define_game,
		"on_play":     on_play,
	})

	// register other stuff
	L.SetField(mod, "name", lua.LString("value"))

	L.SetGlobal("launcher", mod)
	os := L.GetGlobal("os")
	L.SetField(os, "sleep", L.NewFunction(sleep))
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

	fmt.Printf("got id, name, and hero path as strings: %q %q %q\n", id, name, heroPath)

	RegisterGame(Game(id.String(), name.String(), heroPath.String()))

	return 0
}

func on_play(L *lua.LState) int {
	return 0
}
