/* PluginManager.go © Penguin_Spy 2024
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	g "github.com/AllenDang/giu"
	lua "github.com/yuin/gopher-lua"
)

type plugin struct {
	id    string      // the plugin's internal id
	L     *lua.LState // the Lua state of the plugin
	games []*game     // games this plugin has defined. not registered until the plugin successfully finishes loading
}

var plugins = make(map[string]*plugin)
var pluginsDir string

func sleep(L *lua.LState) int {
	duration := L.CheckNumber(1)
	time.Sleep(time.Second * time.Duration(duration))
	return 0
}

func loadLuaLib(L *lua.LState, loader lua.LGFunction) *lua.LTable {
	L.Push(L.NewFunction(loader))
	L.Call(0, 1)
	return L.CheckTable(-1)
}

// creates a new Lua state for a plugin
func newLuaState(id string) *lua.LState {
	L := lua.NewState(lua.Options{SkipOpenLibs: true})
	packageLib := loadLuaLib(L, lua.OpenPackage) // TODO: probably reimplement require and remove package
	packageLib.RawSetString("path", lua.LString(filepath.Join(getPluginDir(id), "?.lua")))
	loadLuaLib(L, lua.OpenBase)
	loadLuaLib(L, lua.OpenTable)
	loadLuaLib(L, lua.OpenString)
	loadLuaLib(L, lua.OpenMath)
	loadLuaLib(L, lua.OpenCoroutine)
	//loadLuaLib(L, lua.OpenIo)
	//loadLuaLib(L, lua.OpenOs)
	//loadLuaLib(L, lua.OpenDebug)
	loadLuaLib(L, OpenLauncher)

	//os := L.GetGlobal("os")
	L.SetGlobal("sleep", L.NewFunction(sleep))

	// used to determine which plugin called our API functions
	L.SetField(L.Get(lua.RegistryIndex), "plugin_id", lua.LString(id))
	return L
}

func loadAllPlugins() {
	pluginsDir = filepath.Join(appDataPath, "plugins")
	if err := os.MkdirAll(pluginsDir, os.ModeDir); err != nil {
		log.Printf("[PluginManager] failed to create plugins directory! - %v\n", err)
		return
	}
	pluginsDirEntries, err := os.ReadDir(pluginsDir)
	if err != nil {
		log.Printf("[PluginManager] failed to read plugins directory! - %v\n", err)
		return
	}

	for _, entry := range pluginsDirEntries {
		if !entry.IsDir() {
			log.Printf("[PluginManager] ignoring non-directory item %q\n", entry)
		} else {
			var id = entry.Name()
			log.Printf("[PluginManager] loading plugin %q\n", id)

			L := newLuaState(id)
			plugin := &plugin{id, L, make([]*game, 0)}
			plugins[id] = plugin
			if err := L.DoFile(filepath.Join(getPluginDir(id), "main.lua")); err != nil {
				log.Printf("[PluginManager] error while loading plugin %q:\n%v", id, err)
				L.Close()
				continue
			} else {
				// once the plugin has successfully loaded, add its games to the list
				for _, game := range plugin.games {
					RegisterGame(game)
				}
			}
		}
	}

	// once all plugins have loaded
	if len(games) > 0 {
		viewingGame = games[0]
	}
	guiState = guiStateGame
	// only update the GUI if it's been created yet
	if g.Context != nil {
		g.Update()
	}
}

func closeAllPlugins() {
	for _, plugin := range plugins {
		if !plugin.L.IsClosed() {
			plugin.L.Close()
		}
	}
}

// gets the absolute path to the directory of the specified plugin
func getPluginDir(id string) string {
	return filepath.Join(pluginsDir, id)
}

func callPluginMethod(plugin *plugin, method string, arg lua.LValue) {
	L := plugin.L

	function, ok := L.GetField(L.Get(lua.RegistryIndex), method).(*lua.LFunction)
	if !ok {
		log.Printf("[PluginManager] could not get plugin %q method %q - %v\n", plugin.id, method, function)
		return
	}

	L.Push(function)
	L.Push(arg)
	// one argument, zero return values, return error normally
	if err := L.PCall(1, 0, nil); err != nil {
		log.Printf("[PluginManager] error while calling plugin %q method %q - %v\n", plugin.id, method, err)
	}
}

func callOnPlay(game *game) {
	plugin := game.plugin
	t := plugin.L.NewTable()
	t.RawSetString("id", lua.LString(game.id))
	t.RawSetString("option", lua.LString("default"))
	callPluginMethod(plugin, "on_play", t)
}
