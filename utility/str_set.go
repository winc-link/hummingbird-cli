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

import "strings"

type StrSet struct {
	mu   RWMutex
	data map[string]struct{}
}

// NewStrSet create and returns a new set, which contains un-repeated items.
// The parameter `safe` is used to specify whether using set in concurrent-safety,
// which is false in default.
func NewStrSet(safe ...bool) *StrSet {
	return &StrSet{
		mu:   Create(safe...),
		data: make(map[string]struct{}),
	}
}

// NewStrSetFrom returns a new set from `items`.
func NewStrSetFrom(items []string, safe ...bool) *StrSet {
	m := make(map[string]struct{})
	for _, v := range items {
		m[v] = struct{}{}
	}
	return &StrSet{
		mu:   Create(safe...),
		data: m,
	}
}

// Iterator iterates the set readonly with given callback function `f`,
// if `f` returns true then continue iterating; or false to stop.
func (set *StrSet) Iterator(f func(v string) bool) {
	set.mu.RLock()
	defer set.mu.RUnlock()
	for k := range set.data {
		if !f(k) {
			break
		}
	}
}

// Add adds one or multiple items to the set.
func (set *StrSet) Add(item ...string) {
	set.mu.Lock()
	if set.data == nil {
		set.data = make(map[string]struct{})
	}
	for _, v := range item {
		set.data[v] = struct{}{}
	}
	set.mu.Unlock()
}

// AddIfNotExist checks whether item exists in the set,
// it adds the item to set and returns true if it does not exist in the set,
// or else it does nothing and returns false.
func (set *StrSet) AddIfNotExist(item string) bool {
	if !set.Contains(item) {
		set.mu.Lock()
		defer set.mu.Unlock()
		if set.data == nil {
			set.data = make(map[string]struct{})
		}
		if _, ok := set.data[item]; !ok {
			set.data[item] = struct{}{}
			return true
		}
	}
	return false
}

// AddIfNotExistFunc checks whether item exists in the set,
// it adds the item to set and returns true if it does not exists in the set and
// function `f` returns true, or else it does nothing and returns false.
//
// Note that, the function `f` is executed without writing lock.
func (set *StrSet) AddIfNotExistFunc(item string, f func() bool) bool {
	if !set.Contains(item) {
		if f() {
			set.mu.Lock()
			defer set.mu.Unlock()
			if set.data == nil {
				set.data = make(map[string]struct{})
			}
			if _, ok := set.data[item]; !ok {
				set.data[item] = struct{}{}
				return true
			}
		}
	}
	return false
}

// AddIfNotExistFuncLock checks whether item exists in the set,
// it adds the item to set and returns true if it does not exists in the set and
// function `f` returns true, or else it does nothing and returns false.
//
// Note that, the function `f` is executed without writing lock.
func (set *StrSet) AddIfNotExistFuncLock(item string, f func() bool) bool {
	if !set.Contains(item) {
		set.mu.Lock()
		defer set.mu.Unlock()
		if set.data == nil {
			set.data = make(map[string]struct{})
		}
		if f() {
			if _, ok := set.data[item]; !ok {
				set.data[item] = struct{}{}
				return true
			}
		}
	}
	return false
}

// Contains checks whether the set contains `item`.
func (set *StrSet) Contains(item string) bool {
	var ok bool
	set.mu.RLock()
	if set.data != nil {
		_, ok = set.data[item]
	}
	set.mu.RUnlock()
	return ok
}

// ContainsI checks whether a value exists in the set with case-insensitively.
// Note that it internally iterates the whole set to do the comparison with case-insensitively.
func (set *StrSet) ContainsI(item string) bool {
	set.mu.RLock()
	defer set.mu.RUnlock()
	for k := range set.data {
		if strings.EqualFold(k, item) {
			return true
		}
	}
	return false
}

// Remove deletes `item` from set.
func (set *StrSet) Remove(item string) {
	set.mu.Lock()
	if set.data != nil {
		delete(set.data, item)
	}
	set.mu.Unlock()
}
