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

package utility

import (
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/genv"
)

const (
	EnvName = "GF_CLI_ALL_YES"
)

// Init initializes the package manually.
func Init() {
	if gcmd.GetOpt("y") != nil {
		genv.MustSet(EnvName, "1")
	}
}

// Check checks whether option allow all yes for command.
func Check() bool {
	return genv.Get(EnvName).String() == "1"
}
