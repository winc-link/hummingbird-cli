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
	"context"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/os/glog"
)

const (
	headerPrintEnvName = "GF_CLI_MLOG_HEADER"
)

var (
	ctx    = context.TODO()
	logger = glog.New()
)

func init() {
	logger.SetStack(false)
	if genv.Get(headerPrintEnvName).String() == "1" {
		logger.SetHeaderPrint(true)
	} else {
		logger.SetHeaderPrint(false)
	}
	if gcmd.GetOpt("debug") != nil || gcmd.GetOpt("gf.debug") != nil {
		logger.SetDebug(true)
	} else {
		logger.SetDebug(false)
	}
}

// SetHeaderPrint enables/disables header printing to stdout.
func SetHeaderPrint(enabled bool) {
	logger.SetHeaderPrint(enabled)
	if enabled {
		_ = genv.Set(headerPrintEnvName, "1")
	} else {
		_ = genv.Set(headerPrintEnvName, "0")
	}
}

func Print(v ...interface{}) {
	logger.Print(ctx, v...)
}

func Printf(format string, v ...interface{}) {
	logger.Printf(ctx, format, v...)
}

func Fatal(v ...interface{}) {
	logger.Fatal(ctx, v...)
}

func Fatalf(format string, v ...interface{}) {
	logger.Fatalf(ctx, format, v...)
}

func Debug(v ...interface{}) {
	logger.Debug(ctx, v...)
}

func Debugf(format string, v ...interface{}) {
	logger.Debugf(ctx, format, v...)
}
