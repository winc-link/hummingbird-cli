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
	"bufio"
	"fmt"
	"os"
)

func Scan(info ...interface{}) string {
	fmt.Print(info...)
	return readline()
}

// Scanf prints `info` to stdout with `format`, reads and returns user input, which stops by '\n'.
func Scanf(format string, info ...interface{}) string {
	fmt.Printf(format, info...)
	return readline()
}

func readline() string {
	var s string
	reader := bufio.NewReader(os.Stdin)
	s, _ = reader.ReadString('\n')
	s = Trim(s)
	return s
}
