# nytcrawler

This program was written by Herbie Ong and Helena Ong to collect specific data
from NYTimes.com. The collected data was then further processed by Helena
using [AntConc](https://www.laurenceanthony.net/software/antconc/) for her
research paper.

This program invokes API calls to the [NYTimes Article Search API](https://developer.nytimes.com/docs/articlesearch-product/1/overview).
It downloads publish date, web URL and headline for specific query term(s) by
year and optionally by month(s). The API response provides a max of 10 results
per page, and max of 201 pages (page 0 to 200) regardless of date range. The
program makes the API call on a per month basis and pages through the results.
It currently cannot handle a response with more than 2010 results.

To use this program, users will need to sign up for an API key. It expects the
user to pass an API key to the -apikey flag when invoked.

Sample invocation:

```sh
$ nytcrawler -apikey=<API_KEY> -query='climate change' -years=2019,2018 output.csv
```

Code was written without much comments due to lack of time. Feel free to use
and modify according to your needs.
