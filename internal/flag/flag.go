// Copyright 2020 The nytcrawler authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package flag

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type uintList []uint

func (f *uintList) Set(value string) error {
	if len(*f) > 0 {
		(*f) = (*f)[0:0]
	}
	nums := strings.Split(value, ",")
	if len(nums) == 1 && strings.TrimSpace(nums[0]) == "" {
		return nil
	}
	for _, s := range nums {
		n, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return fmt.Errorf("error parsing flag value: %v", value)
		}
		*f = append(*f, uint(n))
	}
	return nil
}

func (f *uintList) String() string {
	if len(*f) == 0 {
		return ""
	}
	var list []string
	for _, n := range *f {
		list = append(list, strconv.FormatUint(uint64(n), 10))
	}
	return strings.Join(list, ",")
}

func (f *uintList) Get() interface{} {
	return []uint(*f)
}

func UintList(name string, def []uint, usage string) *[]uint {
	var v []uint
	for _, n := range def {
		v = append(v, n)
	}
	UintListVar(&v, name, usage)
	return &v
}

func UintListVar(v *[]uint, name string, usage string) {
	flag.Var((*uintList)(v), name, usage)
}
