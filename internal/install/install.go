/*******************************************************************************
 * Copyright 2017.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package install

import (
	"fmt"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/spf13/cobra"
	"github.com/winc-link/hummingbird-cli/utility"
	"os"
	"runtime"
	"strings"
)

var CmdInstall = &cobra.Command{
	Use:     "install",
	Example: "./hb install",
	Long:    `install hb CLI.`,
	Short:   `install hb CLI.`,

	Run: run,
}

type serviceInstallAvailablePath struct {
	dirPath   string
	filePath  string
	writable  bool
	installed bool
	IsSelf    bool
}

func Get(key string, def ...interface{}) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return ""
	}
	return v
}

func getGoPathBin() string {
	if goPath := Get(`GOPATH`); goPath != "" {
		return Join(goPath, "bin")
	}
	return ""
}

func run(cmd *cobra.Command, args []string) {
	// Ask where to install.
	paths := getAvailablePaths()
	if len(paths) <= 0 {
		fmt.Printf("no path detected, you can manually install hb by copying the binary to path folder.")
		return
	}
	utility.Printf("I found some installable paths for you(from $PATH): ")
	utility.Printf("  %2s | %8s | %9s | %s", "Id", "Writable", "Installed", "Path")

	// Print all paths status and determine the default selectedID value.
	var (
		selectedID = -1
		newPaths   []serviceInstallAvailablePath
		pathSet    = utility.NewStrSet() // Used for repeated items filtering.
	)
	for _, path := range paths {
		if !pathSet.AddIfNotExist(path.dirPath) {
			continue
		}
		newPaths = append(newPaths, path)
	}
	paths = newPaths
	for id, path := range paths {
		utility.Printf("  %2d | %8t | %9t | %s", id, path.writable, path.installed, path.dirPath)
		if selectedID == -1 {
			// Use the previously installed path as the most priority choice.
			if path.installed {
				selectedID = id
			}
		}
	}

	// If there's no previously installed path, use the first writable path.
	if selectedID == -1 {
		// Order by choosing priority.
		commonPaths := garray.NewStrArrayFrom(g.SliceStr{
			getGoPathBin(),
			`/usr/local/bin`,
			`/usr/bin`,
			`/usr/sbin`,
			`C:\Windows`,
			`C:\Windows\system32`,
			`C:\Go\bin`,
			`C:\Program Files`,
			`C:\Program Files (x86)`,
		})
		// Check the common installation directories.
		commonPaths.Iterator(func(k int, v string) bool {
			for id, aPath := range paths {
				if strings.EqualFold(aPath.dirPath, v) {
					selectedID = id
					return false
				}
			}
			return true
		})
		if selectedID == -1 {
			selectedID = 0
		}
	}

	if utility.Check() {
		// Use the default selectedID.
		utility.Printf("please choose one installation destination [default %d]: %d", selectedID, selectedID)
	} else {
		for {
			// Get input and update selectedID.
			var (
				inputID int
				input   = utility.Scanf("please choose one installation destination [default %d]: ", selectedID)
			)
			if input != "" {
				inputID = gconv.Int(input)
			} else {
				break
			}
			// Check if out of range.
			if inputID >= len(paths) || inputID < 0 {
				utility.Printf("invalid install destination Id: %d", inputID)
				continue
			}
			selectedID = inputID
			break
		}
	}

	// Get selected destination path.
	dstPath := paths[selectedID]

	// Install the new binary.
	utility.Debugf(`copy file from "%s" to "%s"`, gfile.SelfPath(), dstPath.filePath)
	err := CopyFile(SelfPath(), dstPath.filePath)
	if err != nil {
		utility.Printf("install hb binary to '%s' failed: %v", dstPath.dirPath, err)
		utility.Printf("you can manually install hb by copying the binary to folder: %s", dstPath.dirPath)
	} else {
		utility.Printf("hb binary is successfully installed to: %s", dstPath.filePath)
	}
	return

}

// getAvailablePaths returns the installation paths data for the binary.
func getAvailablePaths() []serviceInstallAvailablePath {
	var (
		folderPaths    []serviceInstallAvailablePath
		binaryFileName = "hb" + Ext(SelfPath())
	)
	// $GOPATH/bin
	if goPathBin := getGoPathBin(); goPathBin != "" {
		folderPaths = checkAndAppendToAvailablePath(
			folderPaths, goPathBin, binaryFileName,
		)
	}
	switch runtime.GOOS {
	case "darwin":
		darwinInstallationCheckPaths := []string{"/usr/local/bin"}
		for _, v := range darwinInstallationCheckPaths {
			folderPaths = checkAndAppendToAvailablePath(
				folderPaths, v, binaryFileName,
			)
		}
		fallthrough

	default:
		// Search and find the writable directory path.
		envPath := Get("PATH")
		if envPath == "" {
			envPath = Get("Path")
		}

		if strings.Contains(envPath, ";") {
			// windows.
			for _, v := range utility.SplitAndTrim(envPath, ";") {
				if v == "." {
					continue
				}
				folderPaths = checkAndAppendToAvailablePath(
					folderPaths, v, binaryFileName,
				)
			}
		} else if strings.Contains(envPath, ":") {
			// *nix.
			for _, v := range utility.SplitAndTrim(envPath, ":") {
				if v == "." {
					continue
				}
				folderPaths = checkAndAppendToAvailablePath(
					folderPaths, v, binaryFileName,
				)
			}
		} else if envPath != "" {
			folderPaths = checkAndAppendToAvailablePath(
				folderPaths, envPath, binaryFileName,
			)
		} else {
			folderPaths = checkAndAppendToAvailablePath(
				folderPaths, "/usr/local/bin", binaryFileName,
			)
		}
	}
	return folderPaths
}

func checkAndAppendToAvailablePath(folderPaths []serviceInstallAvailablePath, dirPath string, binaryFileName string) []serviceInstallAvailablePath {
	var (
		filePath  = Join(dirPath, binaryFileName)
		writable  = IsWritable(dirPath)
		installed = Exists(filePath)
		self      = SelfPath() == filePath
	)
	if !writable && !installed {
		return folderPaths
	}
	return append(
		folderPaths,
		serviceInstallAvailablePath{
			dirPath:   dirPath,
			writable:  writable,
			filePath:  filePath,
			installed: installed,
			IsSelf:    self,
		})
}
