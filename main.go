/* main.go © Penguin_Spy 2024
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var appDataPath string

func main() {
	// saves the absolute path to the platform-specific app data directory
	// or the current directory if running from "go run ."
	if strings.Contains(filepath.Dir(os.Args[0]), "go-build") {
		appDataPath = "."
	} else {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			log.Printf("[main] failed to load app data directory! - %v\n", err)
			return
		}
		appDataPath = filepath.Join(userConfigDir, "pylon-launcher")
		if err := os.MkdirAll(userConfigDir, os.ModeDir); err != nil {
			log.Printf("[main] failed to create app data directory! - %v\n", err)
			return
		}
	}

	// save the old log (ignoring any errors, the file might not exist and that's fine)
	os.Rename(filepath.Join(appDataPath, "launcher-current.log"), filepath.Join(appDataPath, "launcher-previous.log"))

	// and open the new log file
	logFile, err := os.Create(filepath.Join(appDataPath, "launcher-current.log"))
	if err != nil {
		log.Printf("[main] failed to open log file (continuing anyways) - %v\n", err)
	}
	defer logFile.Close()
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	log.Printf("[main] using data path %q\n", appDataPath)

	// start the program

	go func() {
		loadAllPlugins()
	}()

	GUI_start()

	closeAllPlugins()
}
