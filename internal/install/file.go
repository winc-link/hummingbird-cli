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
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	// Separator for file system.
	// It here defines the separator as variable
	// to allow it modified by developer if necessary.
	Separator = string(filepath.Separator)

	// DefaultPermOpen is the default perm for file opening.
	DefaultPermOpen = os.FileMode(0666)

	// DefaultPermCopy is the default perm for file/folder copy.
	DefaultPermCopy = os.FileMode(0755)
)

var (
	// The absolute file path for main package.
	// It can be only checked and set once.
	//mainPkgPath = gtype.NewString()

	// selfPath is the current running binary path.
	// As it is most commonly used, it is so defined as an internal package variable.
	selfPath = ""
)

func init() {
	// Initialize internal package variable: selfPath.
	selfPath, _ = exec.LookPath(os.Args[0])
	if selfPath != "" {
		selfPath, _ = filepath.Abs(selfPath)
	}
	if selfPath == "" {
		selfPath, _ = filepath.Abs(os.Args[0])
	}
}

func SelfPath() string {
	return selfPath
}

func Ext(path string) string {
	ext := filepath.Ext(path)
	if p := strings.IndexByte(ext, '?'); p != -1 {
		ext = ext[0:p]
	}
	return ext
}

func OpenWithFlagPerm(path string, flag int, perm os.FileMode) (*os.File, error) {
	file, err := OpenFile(path, flag, perm)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func OpenFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	file, err := os.OpenFile(path, flag, perm)
	if err != nil {
		err = fmt.Errorf("os.OpenFile failed with name %s, flag %d, perm %d ", path, flag, perm)
	}
	return file, err
}

// Join joins string array paths with file separator of current system.
func Join(paths ...string) string {
	var s string
	for _, path := range paths {
		if s != "" {
			s += Separator
		}
		s += TrimRight(path, Separator)
	}
	return s
}

var (
	// DefaultTrimChars are the characters which are stripped by Trim* functions in default.
	DefaultTrimChars = string([]byte{
		'\t', // Tab.
		'\v', // Vertical tab.
		'\n', // New line (line feed).
		'\r', // Carriage return.
		'\f', // New page.
		' ',  // Ordinary space.
		0x00, // NUL-byte.
		0x85, // Delete.
		0xA0, // Non-breaking space.
	})
)

func TrimRight(str string, characterMask ...string) string {
	trimChars := DefaultTrimChars
	if len(characterMask) > 0 {
		trimChars += characterMask[0]
	}
	return strings.TrimRight(str, trimChars)
}

// Exists checks whether given `path` exist.
func Exists(path string) bool {
	if stat, err := os.Stat(path); stat != nil && !os.IsNotExist(err) {
		return true
	}
	return false
}

// IsDir checks whether given `path` a directory.
// Note that it returns false if the `path` does not exist.
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsWritable(path string) bool {
	result := true
	if IsDir(path) {
		// If it's a directory, create a temporary file to test whether it's writable.
		tmpFile := strings.TrimRight(path, Separator) + Separator + strconv.FormatInt(time.Now().UnixNano(), 10)
		if f, err := Create(tmpFile); err != nil || !Exists(tmpFile) {
			result = false
		} else {
			_ = f.Close()
			_ = Remove(tmpFile)
		}
	} else {
		// If it's a file, check if it can open it.
		file, err := os.OpenFile(path, os.O_WRONLY, DefaultPermOpen)
		if err != nil {
			result = false
		}
		_ = file.Close()
	}
	return result
}

func Create(path string) (*os.File, error) {
	dir := Dir(path)
	if !Exists(dir) {
		if err := Mkdir(dir); err != nil {
			return nil, err
		}
	}
	file, err := os.Create(path)
	if err != nil {
		err = fmt.Errorf("os.Create failed for name %s", path)
	}
	return file, err
}

// Open opens file/directory READONLY.
func Open(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		err = fmt.Errorf("os.Open failed for name %s", path)

	}
	return file, err
}

func Dir(path string) string {
	if path == "." {
		return filepath.Dir(RealPath(path))
	}
	return filepath.Dir(path)
}

func RealPath(path string) string {
	p, err := filepath.Abs(path)
	if err != nil {
		return ""
	}
	if !Exists(p) {
		return ""
	}
	return p
}

func Mkdir(path string) (err error) {
	if err = os.MkdirAll(path, os.ModePerm); err != nil {
		err = fmt.Errorf("os.MkdirAll failed for path %s with perm %d", path, os.ModePerm)
		return err
	}
	return nil
}

func Remove(path string) (err error) {
	// It does nothing if `path` is empty.
	if path == "" {
		return nil
	}
	if err = os.RemoveAll(path); err != nil {
		err = fmt.Errorf("os.RemoveAll failed for path %s", path)
	}
	return
}

func CopyFile(src, dst string, option ...CopyOption) (err error) {
	var usedOption = getCopyOption(option...)
	if src == "" {
		return errors.New("source file cannot be empty")
	}
	if dst == "" {
		return errors.New("destination file cannot be empty")

	}
	// If src and dst are the same path, it does nothing.
	if src == dst {
		return nil
	}
	// file state check.
	srcStat, srcStatErr := os.Stat(src)
	if srcStatErr != nil {
		if os.IsNotExist(srcStatErr) {
			return errors.New(fmt.Sprintf("the src path %s does not exist", src))
		}
		return errors.New(fmt.Sprintf("call os.Stat on %s failed", src))

	}
	dstStat, dstStatErr := os.Stat(dst)
	if dstStatErr != nil && !os.IsNotExist(dstStatErr) {
		return fmt.Errorf("call os.Stat on %s failed", dst)
	}
	if !srcStat.IsDir() && dstStat != nil && dstStat.IsDir() {
		return fmt.Errorf("CopyFile failed: the src path %s is file, but the dst path %s is folder", src, dst)

	}
	// copy file logic.
	var inFile *os.File
	inFile, err = Open(src)
	if err != nil {
		return
	}
	defer func() {
		if e := inFile.Close(); e != nil {
			err = fmt.Errorf("file close failed for %s", src)
		}
	}()
	var outFile *os.File
	outFile, err = Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := outFile.Close(); e != nil {
			err = fmt.Errorf("file close failed for %s", dst)

		}
	}()
	if _, err = io.Copy(outFile, inFile); err != nil {
		err = fmt.Errorf("io.Copy failed from %s to %s", src, dst)
		return
	}
	if usedOption.Sync {
		if err = outFile.Sync(); err != nil {
			err = fmt.Errorf("file sync failed for file %s", dst)
			return
		}
	}
	if usedOption.PreserveMode {
		usedOption.Mode = srcStat.Mode().Perm()
	}
	if err = Chmod(dst, usedOption.Mode); err != nil {
		return
	}
	return
}

type CopyOption struct {
	// Auto call file sync after source file content copied to target file.
	Sync bool

	// Preserve the mode of the original file to the target file.
	// If true, the Mode attribute will make no sense.
	PreserveMode bool

	// Destination created file mode.
	// The default file mode is DefaultPermCopy if PreserveMode is false.
	Mode os.FileMode
}

func getCopyOption(option ...CopyOption) CopyOption {
	var usedOption CopyOption
	if len(option) > 0 {
		usedOption = option[0]
	}
	if usedOption.Mode == 0 {
		usedOption.Mode = DefaultPermCopy
	}
	return usedOption
}

func Chmod(path string, mode os.FileMode) (err error) {
	err = os.Chmod(path, mode)
	if err != nil {
		err = fmt.Errorf("os.Chmod failed with path %s and mode %s", path, mode)
	}
	return
}
