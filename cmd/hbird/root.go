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

package hbird

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/winc-link/hummingbird-cli/config"
	"github.com/winc-link/hummingbird-cli/internal/install"
	"github.com/winc-link/hummingbird-cli/internal/new"
)

var CmdRoot = &cobra.Command{
	Use:     "hb",
	Example: "hb new demo-driver",
	Short:   config.LogoContent,
	Version: config.LogoContent + "\n" + fmt.Sprintf("Hummingbird %s - Copyright (c) 2023 hb \nReleased under the MIT License.\n", config.Version),
}

func init() {
	CmdRoot.AddCommand(new.CmdNew)
	CmdRoot.AddCommand(install.CmdInstall)

}

// Execute executes the root command.
func Execute() error {
	return CmdRoot.Execute()
}
