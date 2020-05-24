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
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUintList(t *testing.T) {
	tests := []struct {
		in   string
		want *uintList
	}{
		{
			in:   "",
			want: &uintList{},
		},
		{
			in:   "  ",
			want: &uintList{},
		},
		{
			in:   "0",
			want: &uintList{0},
		},
		{
			in:   "0,101,42,47",
			want: &uintList{0, 101, 42, 47},
		},
	}

	for _, tc := range tests {
		t.Run("input:"+tc.in, func(t *testing.T) {
			f := &uintList{}
			if err := f.Set(tc.in); err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.want, f); diff != "" {
				t.Errorf("flag (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUintListErr(t *testing.T) {
	tests := []string{
		"hello",
		"-47",
		"101.1",
		"0,,42,47",
		"0,abc,42",
		"0,  101 ,42",
	}

	for _, tc := range tests {
		t.Run("input:"+tc, func(t *testing.T) {
			f := &uintList{}
			if err := f.Set(tc); err == nil {
				t.Error("expecting error, got nil")
			}
		})
	}
}
