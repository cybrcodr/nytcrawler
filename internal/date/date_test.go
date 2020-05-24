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

package date

import (
	"fmt"
	"testing"
)

func TestIsLeapYear(t *testing.T) {
	for _, i := range []uint{1600, 2000, 2012, 2016, 2020} {
		if !isLeapYear(i) {
			t.Errorf("isLeapYear(%d) got false, want true", i)
		}
	}

	for _, i := range []uint{1601, 2002, 2015, 2017, 2018} {
		if isLeapYear(i) {
			t.Errorf("isLeapYear(%d) got true, want false", i)
		}
	}
}

func TestMonthRange(t *testing.T) {
	tests := []struct {
		year  uint
		month uint
		want  Range
	}{
		{2015, 1, Range{"20150101", "20150131"}},
		{2018, 12, Range{"20181201", "20181231"}},
		{2015, 2, Range{"20150201", "20150228"}},
		{2016, 2, Range{"20160201", "20160229"}},
		{2020, 2, Range{"20200201", "20200229"}},
		{2020, 4, Range{"20200401", "20200430"}},
		{2021, 6, Range{"20210601", "20210630"}},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%d-%02d", tc.year, tc.month), func(t *testing.T) {
			got := MonthRange(tc.year, tc.month)
			if got != tc.want {
				t.Errorf("MonthRange(%d, %d) got %+v, want %+v", tc.year, tc.month, got, tc.want)
			}
		})
	}
}
