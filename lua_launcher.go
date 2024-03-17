package main

import lua "github.com/yuin/gopher-lua"

func launcher_open_lua_lib(L *lua.LState) int {
	// register functions to the table
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"define_game": define_game,
		"on_play":     on_play,
	})

	// register other stuff
	L.SetField(mod, "name", lua.LString("value"))

	L.SetGlobal("launcher", mod)
	return 1
}

func define_game(L *lua.LState) int {
	return 0
}

func on_play(L *lua.LState) int {
	return 0
}
