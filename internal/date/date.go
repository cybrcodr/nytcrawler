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

import "fmt"

// Range contains the start and end dates in string with formats YYYYMMDD.
type Range struct {
	Start string
	End   string
}

func (r Range) String() string {
	return r.Start + "-" + r.End
}

// MonthRange returns a Range instance for the given month of the year.
func MonthRange(year, month uint) Range {
	if month < 1 || month > 12 {
		panic(fmt.Sprintf("MonthRange() given invalid month %d", month))
	}

	days := daysInMonth[month]
	if month == 2 && isLeapYear(year) {
		days++
	}
	return Range{
		Start: fmt.Sprintf("%04d%02d01", year, month),
		End:   fmt.Sprintf("%04d%02d%02d", year, month, days),
	}
}

var daysInMonth = map[uint]uint{
	1:  31,
	2:  28, // non-leap year.
	3:  31,
	4:  30,
	5:  31,
	6:  30,
	7:  31,
	8:  31,
	9:  30,
	10: 31,
	11: 30,
	12: 31,
}

func isLeapYear(year uint) bool {
	if year%4 != 0 || (year%100 == 0 && year%400 != 0) {
		return false
	}
	return true
}
