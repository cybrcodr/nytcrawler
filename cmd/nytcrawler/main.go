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

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cybrcodr/nytcrawler/internal/api"
	"github.com/cybrcodr/nytcrawler/internal/date"
	iflag "github.com/cybrcodr/nytcrawler/internal/flag"
)

var (
	years  = iflag.UintList("years", nil, "years to search on")
	months = iflag.UintList("months",
		[]uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		"months to search on")
	query  = flag.String("query", "", "query term")
	apiKey = flag.String("apikey", "", "API key")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	apiKey := checkAPIKey()
	query := checkQueryTerm()
	years := checkYears()
	months := checkMonths()
	file := checkFile()
	defer file.Close()
	w := csv.NewWriter(file)

	log.Printf("writing results to %s", flag.Arg(0))
	start := time.Now()
loop:
	for _, y := range years {
		for _, m := range months {
			d := date.MonthRange(y, m)
			err := api.FetchArticles(apiKey, query, d, func(results []api.Result) error {
				for _, res := range results {
					w.Write([]string{
						res.PubDate,
						res.WebURL,
						res.Headline,
					})
				}
				w.Flush()
				return w.Error()
			})
			if err != nil {
				log.Printf("%s: fetch error: %v\n", d, err)
				break loop
			}
		}
	}

	api.PrintStats()
	log.Printf("Total time of %v", time.Now().Sub(start))
	log.Println("That's all folks!!!")
}

func checkAPIKey() string {
	ret := strings.TrimSpace(*apiKey)
	if len(ret) == 0 {
		fmt.Fprintf(os.Stderr, "missing -apikey value\n")
		usage()
		os.Exit(1)
	}

	return ret
}

func checkQueryTerm() string {
	ret := strings.TrimSpace(*query)
	if len(ret) == 0 {
		fmt.Fprintf(os.Stderr, "missing -query value\n")
		usage()
		os.Exit(1)
	}

	return ret
}

func checkYears() []uint {
	if len(*years) == 0 {
		fmt.Fprintf(os.Stderr, "missing -years value\n")
		usage()
		os.Exit(1)
	}

	for _, y := range *years {
		if y < 1900 || y > 2021 {
			fmt.Fprintf(os.Stderr, "year not within range: %d\n", y)
			os.Exit(1)
		}
	}
	return *years
}

func checkMonths() []uint {
	if len(*months) == 0 {
		fmt.Fprintf(os.Stderr, "missing -months value\n")
		usage()
		os.Exit(1)
	}

	for _, m := range *months {
		if m < 1 || m > 12 {
			fmt.Fprintf(os.Stderr, "month not within range: %d\n", m)
			os.Exit(1)
		}
	}
	return *months
}

func checkFile() *os.File {
	if len(flag.Args()) != 1 {
		usage()
		os.Exit(1)
	}
	file := flag.Arg(0)
	out, err := os.Create(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating file %s: %v", file, err)
		os.Exit(1)
	}
	return out
}

func usage() {
	fmt.Fprintf(os.Stderr,
		"Usage: %v -apikey=<API-key> -query=<query-term> -years=<Y1,Y2,...> -months=<1,2,...> <file>\n\n",
		filepath.Base(os.Args[0]))

	flag.PrintDefaults()
}
