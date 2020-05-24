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

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/cybrcodr/nytcrawler/internal/date"
)

const (
	apiURL = `https://api.nytimes.com/svc/search/v2/articlesearch.json`
	fields = `headline,pub_date,web_url,uri`

	docsPerPage = 10
	maxPage     = 200
)

var queryFilters = [...]string{
	`type_of_material:("News")`,
	`document_type:("article")`,
}

var searchURL string

func init() {
	var sb strings.Builder
	sb.WriteString(apiURL)
	sb.WriteString(`?`)
	sb.WriteString(`&sort=oldest`)
	sb.WriteString(`&facet_filter=true`)
	sb.WriteString(`&fl=` + url.QueryEscape(fields))
	for _, filter := range queryFilters {
		sb.WriteString(`&fq=` + url.QueryEscape(filter))
	}

	searchURL = sb.String()
}

type Result struct {
	URI      string
	PubDate  string
	Headline string
	WebURL   string
}

func FetchArticles(apiKey string, term string, d date.Range, output func([]Result) error) error {
	resp, err := fetch(apiKey, term, d, 0)
	if err != nil {
		return err
	}

	hits := resp.Response.Meta.Hits
	if hits == 0 {
		return nil
	}

	results := make([]Result, 0, docsPerPage)
	for _, doc := range resp.Response.Docs {
		results = append(results, Result{
			PubDate:  doc.PubDate,
			Headline: removeNewLines(doc.Headline.Main),
			WebURL:   doc.WebURL,
		})
		stats.hits++
	}
	if err := output(results); err != nil {
		return err
	}

	if hits <= docsPerPage {
		return nil
	}

	// Compute for how many more pages to make requests.
	// Can make requests from 1 to 200 at most.
	pages := hits / docsPerPage
	if hits%docsPerPage == 0 {
		pages--
	}
	if pages > maxPage {
		return fmt.Errorf("hits will exceed max pages to fetch: %d", hits)
	}

	for p := 1; p <= pages; p++ {
		resp, err := fetch(apiKey, term, d, uint(p))
		if err != nil {
			return err
		}

		results = results[:0]
		for _, doc := range resp.Response.Docs {
			results = append(results, Result{
				PubDate:  doc.PubDate,
				Headline: removeNewLines(doc.Headline.Main),
				WebURL:   doc.WebURL,
			})
			stats.hits++
		}
		if err := output(results); err != nil {
			return err
		}
	}

	return nil
}

func removeNewLines(s string) string {
	list := strings.Split(s, "\n")
	return strings.Join(list, " ")
}

const waitTime = 6 * time.Second

func fetch(apiKey string, term string, d date.Range, page uint) (*apiResponse, error) {
	if dur := time.Now().Sub(stats.lastRequest); dur <= waitTime {
		// log.Printf("sleeping %v", dur)
		time.Sleep(waitTime - dur)
	}

	url := searchURL +
		`&api-key=` + url.QueryEscape(apiKey) +
		`&q=` + url.QueryEscape(term) +
		`&begin_date=` + d.Start +
		`&end_date=` + d.End +
		`&page=` + strconv.FormatUint(uint64(page), 10)
	// log.Printf("[%v,page=%d]: URL %v", d, page, url)

	stats.requests++
	stats.lastRequest = time.Now()

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %v: %v", resp.StatusCode, resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	resp.Body.Close()

	apiResp := &apiResponse{}
	if err := json.Unmarshal(b, apiResp); err != nil {
		return nil, fmt.Errorf("error reading JSON response: %v", err)
	}
	// log.Printf("[%v,page=%d]: response %+v", d, 0, apiResp)
	meta := apiResp.Response.Meta
	log.Printf("%v page=%d status=%v hits=%v offset=%v",
		d, page, apiResp.Status, meta.Hits, meta.Offset)

	if apiResp.Status != "OK" {
		return nil, fmt.Errorf("API status: %v", apiResp.Status)
	}

	return apiResp, nil
}

type apiResponse struct {
	Status   string   `json:"status"`
	Response response `json:"response"`
}

type response struct {
	Docs []doc `json:"docs"`
	Meta meta  `json:"meta"`
}

type meta struct {
	Hits   int `json:"hits"`
	Offset int `json:"offset"`
}

type doc struct {
	WebURL   string   `json:"web_url"`
	PubDate  string   `json:"pub_date"`
	Headline headline `json:"headline"`
}

type headline struct {
	Main string `json:"main"`
}

var stats = stats_{}

type stats_ struct {
	requests    uint
	hits        uint
	lastRequest time.Time
}

func PrintStats() {
	log.Printf("number of requests: %d", stats.requests)
	log.Printf("number of hits: %d", stats.hits)
	log.Printf("last request at %v", stats.lastRequest)
}
