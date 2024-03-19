/* PluginManager.go © Penguin_Spy 2024
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	g "github.com/AllenDang/giu"
	lua "github.com/yuin/gopher-lua"
)

const pluginsDir = "plugins"

type plugin struct {
	id string      // the plugin's internal id
	L  *lua.LState // the Lua state of the plugin
}

var plugins = make(map[string]*plugin)

func sleep(L *lua.LState) int {
	duration := L.CheckNumber(1)
	time.Sleep(time.Second * time.Duration(duration))
	return 0
}

func newLuaState(id string) *lua.LState {
	L := lua.NewState()

	launcher_openLuaLib(L)

	os := L.GetGlobal("os")
	L.SetField(os, "sleep", L.NewFunction(sleep))

	L.G.Registry.RawSetString("plugin_id", lua.LString(id))
	return L
}

func loadAllPlugins() {
	dir, err := os.ReadDir(pluginsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[PluginManager] failed to load plugin directory! - %v\n", err)
		return
	}

	for _, entry := range dir {
		if !entry.IsDir() {
			fmt.Printf("[PluginManager] ignoring non-directory item %v\n", entry)
		} else {
			var id = entry.Name()
			var path = filepath.Join(pluginsDir, id, "main.lua")
			fmt.Printf("[PluginManager] loading plugin %q\n", id)

			L := newLuaState(id)
			plugins[id] = &plugin{id, L}
			if err := L.DoFile(path); err != nil {
				fmt.Fprintf(os.Stderr, "[PluginManager] failed to load plugin %v", path)
				L.Close()
				continue
			}
		}
	}

	// once all plugins have loaded
	viewingGame = games[0] // TODO: this crashes if no games were registered
	guiState = guiStateGame
	g.Update()
}

func closeAllPlugins() {
	for _, plugin := range plugins {
		plugin.L.Close()
	}
}

func callPluginMethod(plugin *plugin, method string, arg lua.LValue) {
	L := plugin.L

	function, ok := L.G.Registry.RawGetString(method).(*lua.LFunction)
	if !ok {
		fmt.Fprintf(os.Stderr, "[PluginManager] could not get plugin %q method %q - %v\n", plugin.id, method, function)
		return
	}

	L.Push(function)
	L.Push(arg)
	// one argument, zero return values, return error normally
	if err := L.PCall(1, 0, nil); err != nil {
		fmt.Fprintf(os.Stderr, "[PluginManager] error while calling plugin %q method %q - %v\n", plugin.id, method, err)
	}
}

func callOnPlay(game *game) {
	plugin := plugins[game.pluginId]
	t := plugin.L.NewTable()
	t.RawSetString("id", lua.LString(game.id))
	t.RawSetString("option", lua.LString("default"))
	callPluginMethod(plugin, "on_play", t)
}
