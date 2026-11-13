<p align="center">
  <img src="docs/brand/rajaraja-icon-256.png" width="96" alt="Rajaraja logo">
</p>

<h1 align="center">Rajaraja - Web Crawler</h1>

<p align="center">Survey a whole website from one command. Every page once, every link recorded.</p>

<p align="center">A fast, concurrent Go crawler that maps a site's internal links into one clean JSON report.</p>

<p align="center"><b>Website:</b> <a href="https://rajaraja.surajvijayan.dev">https://rajaraja.surajvijayan.dev</a></p>

```text
$ rajaraja-web-crawler https://learnwebscraping.dev/practice/ecommerce/ 5 50
starting crawl of: https://learnwebscraping.dev/practice/ecommerce/
starting crawl of: https://learnwebscraping.dev/practice/ecommerce/products/sunspire-relic-blade-fan-1006/
starting crawl of: https://www.boot.dev
SKIPPED - not part of base URL
...
SKIPPED - page already visited
report saved to file: report.json
```

<p align="center"><b>45 pages, 394 links and 36 images mapped in 2.4 seconds.</b></p>

Rajaraja is a small command-line crawler written in Go. Give it a start URL and it walks every link that stays on the same site, fetches each HTML page exactly once, and writes `report.json`: for every page, its URL, main heading, first paragraph, every outgoing link and every image. That file is the raw material for internal linking work, content inventories and technical SEO checks, and it drops straight into `jq`, a spreadsheet or a graph tool.

It is one static binary with one third-party library, no config files and no runtime to install.

[Rajaraja](https://en.wikipedia.org/wiki/Rajaraja_I) (pronounced RAH-jah-RAH-jah) is named after Rajaraja Chola I, the Chola emperor whose survey in the year 1000 measured and recorded the land of his whole empire. This tool does the same for a website: it visits every page it can reach, records what it finds, and hands you the map.

## Contents

- [Why I built this](#why-i-built-this)
- [Quick start](#quick-start)
- [Usage](#usage)
- [The report](#the-report)
- [Turning the report into SEO answers](#turning-the-report-into-seo-answers)
- [How it crawls](#how-it-crawls)
- [What it does not do (yet)](#what-it-does-not-do-yet)
- [Crawl responsibly](#crawl-responsibly)
- [Performance](#performance)
- [The name](#the-name)
- [For developers](#for-developers)
- [The website](#the-website)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgements](#acknowledgements)

## Why I built this

I built Rajaraja to learn Go concurrency properly. I'm building [double-distributed-kv](https://github.com/its-me-sv/double-distributed-kv), a key-value store in Go, and before taking it further I wanted hands-on command of every concurrency tool Go gives you: goroutines, channels, mutexes, WaitGroups and semaphores.

A web crawler is the best practice ground I know for that, because it needs all of them at once:

- **Goroutines** that spawn more goroutines, since every page discovers new work.
- **A buffered channel used as a semaphore**, so no more than `maxConcurrency` pages are fetched at a time.
- **A WaitGroup** to know when the last worker has finished, with no queue to close.
- **A map behind a mutex**, read and written by dozens of workers at the same moment.

The bugs that matter only show up under concurrency: two workers fetching the same page, or racing past the page limit together. Rajaraja solves both by checking the limit and claiming the URL inside one locked section, which is exactly the kind of shared state problem a key-value store has to get right.

It also turned out useful on its own. Almost every SEO question about a site starts with the same one: what pages does it have, and how do they link to each other? Rajaraja answers that with one command and one file you can query however you like.

[Suraj Vijayan](https://www.linkedin.com/in/suraj-vijay-32b679224/)

## Quick start

You need [Go](https://go.dev/dl/) 1.27.1 or newer.

```bash
git clone https://github.com/its-me-sv/rajaraja-web-crawler.git
cd rajaraja-web-crawler
go build -o rajaraja-web-crawler .
./rajaraja-web-crawler https://learnwebscraping.dev/practice/ecommerce/ 5 50
```

Or install it straight into your Go bin folder:

```bash
go install github.com/its-me-sv/rajaraja-web-crawler@latest
```

On Windows the binary is `rajaraja-web-crawler.exe`. The example target is a practice site built for scraping, so it is safe to try.

## Usage

```text
rajaraja-web-crawler <baseURL> <maxConcurrency> <maxPages>
```

All three arguments are required, in this order. Run it with `-h` or `--help` to print the usage text; any other number of arguments prints it too and exits with code 1.

| Argument         | What it does                                                               | Rule                                      |
| ---------------- | -------------------------------------------------------------------------- | ----------------------------------------- |
| `baseURL`        | Where the crawl starts. Its hostname decides which links count as internal | A full URL, such as `https://example.com` |
| `maxConcurrency` | How many pages can be fetched at the same time                             | Whole number, 1 or more                   |
| `maxPages`       | The most pages the report will hold                                        | Whole number, 1 or more                   |

Examples:

```bash
# a small site, gently: one request at a time, up to 100 pages
rajaraja-web-crawler https://example.com 1 100

# faster: five requests at once, up to 500 pages
rajaraja-web-crawler https://example.com 5 500
```

While it runs, the terminal shows a progress log:

| Line                                   | Meaning                                   |
| -------------------------------------- | ----------------------------------------- |
| `starting crawl of: <url>`             | A worker picked up a link                 |
| `SKIPPED - not part of base URL`       | The link points to another host           |
| `SKIPPED - page already visited`       | That page was already fetched or claimed  |
| `SKIPPED - failed to get page content` | Error status, timeout or not an HTML page |
| `SKIPPED - invalid page URL`           | The link could not be parsed              |
| `STOPPED - limit reached`              | `maxPages` is full                        |
| `report saved to file: report.json`    | Done                                      |

Workers run at the same time, so log lines interleave. The program exits with code 0 when the report is saved, and 1 if the arguments are invalid or the report cannot be written. Pages that fail during the crawl are skipped and do not change the exit code.

## The report

Every run writes `report.json` in the current folder, replacing any earlier one. It is a JSON array with one object per page, sorted by URL:

```json
{
  "url": "https://learnwebscraping.dev/practice/ecommerce/categories/enchanted-rings/",
  "heading": "Enchanted Rings",
  "first_paragraph": "Browse the Enchanted Rings collection with product pages, prices, rarity labels, and stock information.",
  "outgoing_links": [
    "https://learnwebscraping.dev/practice/ecommerce/",
    "https://learnwebscraping.dev/practice/ecommerce/categories/",
    "..."
  ],
  "image_urls": []
}
```

| Field             | What it holds                                                                                                             |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------- |
| `url`             | The page URL as it was linked                                                                                             |
| `heading`         | Text of the first `h1`, or the first `h2` if there is no `h1`. Empty if neither exists                                    |
| `first_paragraph` | Text of the first `p` inside `<main>`, or the first `p` on the page                                                       |
| `outgoing_links`  | Every `<a href>` on the page as an absolute `http` or `https` URL, internal and external, in page order, repeats included |
| `image_urls`      | Every `<img src>` as an absolute `http` or `https` URL                                                                    |

Only pages that were fetched successfully as HTML appear in the report.

## Turning the report into SEO answers

The report is deliberately raw, so you can ask it your own questions. A few recipes with [`jq`](https://jqlang.org):

```bash
# how many pages were found
jq length report.json

# pages with no h1 or h2
jq -r '.[] | select(.heading == "") | .url' report.json

# most linked-to URLs, counted once per linking page
jq -r '[.[] | .outgoing_links | unique[]] | group_by(.) | map({url: .[0], pages: length}) | sort_by(-.pages) | .[:10][] | "\(.pages)\t\(.url)"' report.json

# every external link the site points to
jq -r --arg site "https://example.com" '[.[].outgoing_links[] | select(startswith($site) | not)] | unique[]' report.json

# every image used across the site
jq -r '[.[].image_urls[]] | unique[]' report.json

# an edge list for Gephi, Cytoscape or a spreadsheet
jq -r '.[] | .url as $from | .outgoing_links | unique[] | [$from, .] | @csv' report.json > edges.csv
```

The edge list is the internal link graph. Load it into a graph tool to see hub pages, dead ends and how far content sits from the home page.

## How it crawls

**What counts as internal.** A link is followed when its hostname exactly matches the start URL's hostname. `http` and `https` versions of a page count as the same page. Subdomains, and `www` versus the bare domain, count as different sites, so start on the host the site actually uses. The start path does not limit the crawl: starting at `/blog/` still follows links to anywhere on the same host.

**When two URLs are the same page.** Before fetching, each URL is turned into a key:

- the scheme is dropped, so `http` and `https` match
- the host is lowercased, the path keeps its case (`/A` and `/a` are different pages)
- default ports (`:80`, `:443`) are dropped, other ports are kept
- trailing slashes and `#fragments` are removed
- query strings are kept as they are, because they can change the page

So `https://Example.com/blog/`, `http://example.com/blog` and `https://example.com/blog#top` are one page, fetched once.

**Never twice, never over the limit.** All workers share one map of pages behind a lock. A worker checks the page limit and claims a URL in the same step, so two workers can never fetch the same page and the report never holds more than `maxPages`. If a fetch fails, the claim is released.

**Concurrency.** Each discovered link gets its own goroutine, and a semaphore sized to `maxConcurrency` caps how many are working at once. All workers share one HTTP client that keeps connections open, so pages on the same host do not pay for a new TLS handshake each time.

**Requests.** Every request sends the user agent `RajarajaWebCrawler/1.0` and has a 7 second timeout. Redirects are followed (up to 10 hops). Pages with a status of 400 or above, and responses whose `Content-Type` is not `text/html`, are skipped.

**Links and images.** Relative links are resolved against the page URL. `mailto:`, `tel:`, `javascript:` and other non-web schemes are dropped.

## What it does not do (yet)

Rajaraja maps what a site links to. It does not judge it. To keep expectations honest:

- **No JavaScript rendering.** It reads the HTML the server sends. Sites that build their pages in the browser (React, Vue or similar single-page apps without server rendering) show up as a single page with no links.
- **No robots.txt support**, no crawl delay and no retries. See [Crawl responsibly](#crawl-responsibly).
- **Broken links and redirects are not reported.** Failed pages are left out of the report, and a redirected URL is stored with the content of the page it redirected to.
- **No depth limit.** `maxPages` is the only bound, so sites with endless query string variants (calendars, faceted filters) can fill the limit.
- **Ignored markup:** `<title>`, meta description, canonical, meta robots, `rel="nofollow"` and `<base href>`.
- **Capped crawls vary.** When a site has more pages than `maxPages`, which pages make the cut can change from run to run. Uncapped crawls are identical every time.
- **The report path is fixed** to `./report.json`, and stopping with Ctrl+C writes nothing.

## Crawl responsibly

Because Rajaraja does not read robots.txt or pause between requests, you are in charge of being polite:

- Crawl sites you own, or sites whose owners allow it.
- Start with a low `maxConcurrency` (1 or 2) on sites you do not run.
- Check the site's `/robots.txt` yourself before a large crawl.

## Performance

Measured on a home connection in India against the 45-page practice site above:

| Command                            | Pages |   Time |
| ---------------------------------- | ----: | -----: |
| `rajaraja-web-crawler <site> 1 50` |    45 | 8.56 s |
| `rajaraja-web-crawler <site> 5 50` |    45 | 2.40 s |

Five workers were about 3.6 times faster than one, and both reports were byte for byte identical.

## The name

Rajaraja Chola I, in Tamil **ராஜராஜ சோழன்** (Rajaraja Chozhan), ruled the Chola empire from 985 to 1014 and built the great Rajarajeshwaram temple at Thanjavur. Around the year 1000 he began a survey and assessment of the land across his empire, and the land was reorganised into revenue divisions as a result. The work was led by an officer called Kuravan Ulagalandan, a title read as "he who measured the world", and the records were cut into the stone walls of temples. Later rulers updated the survey from time to time.

| Then                                                                                      | Now                                        |
| ----------------------------------------------------------------------------------------- | ------------------------------------------ |
| **வளநாடு** (valanadu, "fertile land"), the divisions the survey drew                      | The site boundary: one hostname            |
| **உலகளந்தான்** (Ulagalandan, "he who measured the world"), the officer who ran the survey | The workers that walk the site in parallel |
| Every plot measured once and recorded                                                     | Every page fetched once and recorded       |
| Records cut into temple walls                                                             | `report.json`                              |
| The survey revisited by later rulers                                                      | Run it again whenever the site changes     |

Sources: [Rajaraja I](https://en.wikipedia.org/wiki/Rajaraja_I), [Chola government](https://en.wikipedia.org/wiki/Chola_government), [Middle Chola Temples on the Rajaraja survey](https://de.wisdomlib.org/south-asia/book/middle-chola-temples/d/doc210447.html).

## For developers

### Requirements

Go 1.27.1 or newer (set in `go.mod`).

### Build and test

```bash
go build -o rajaraja-web-crawler .
go test ./...
go vet ./...
```

The test suite has 163 tests and subtests covering URL normalisation and every HTML extractor, including malformed HTML, HTML entities, comments, protocol-relative links and `../` paths that climb past the root.

### Dependencies

| Module                                                            | Why                                                                                                |
| ----------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| [goquery](https://github.com/PuerkitoBio/goquery)                 | CSS selectors over parsed HTML, so each extractor is a few lines (`a[href]`, `img[src]`, `main p`) |
| [cascadia](https://github.com/andybalholm/cascadia)               | Selector engine behind goquery (indirect)                                                          |
| [golang.org/x/net/html](https://pkg.go.dev/golang.org/x/net/html) | Standards-based HTML parser behind goquery (indirect)                                              |

Everything else is the Go standard library.

### Project layout

```
main.go                     entry point: read arguments, crawl, write report.json
internal/
  crawler/
    crawler.go              config, shared page map and lock, shared HTTP client
    crawl.go                Crawl() and the worker goroutine
  extractor/
    extractors.go           PageData and ExtractPageData
    get_heading_from_html.go, get_first_paragraph_from_html.go,
    get_urls_from_html.go, get_images_from_html.go, plus tests
  utils/
    parse_args.go           positional arguments
    get_html.go             HTTP GET with user agent, status and content type checks
    normalize_url.go        the page key used to avoid repeat fetches, plus tests
    json_report.go          sorted, indented JSON output
site/                       the project website (static HTML, CSS and JS, no build step)
.github/                    Pages deploy workflow and issue templates
```

### Design choices

| Choice                                                                              | Instead of                                | Why                                                                                                                                                                                 |
| ----------------------------------------------------------------------------------- | ----------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| A goroutine per link, gated by a semaphore and a WaitGroup                          | A fixed worker pool reading from a queue  | Fits a crawl where every page discovers more work. The crawl is done when the WaitGroup reaches zero, with no queue closing logic. The trade is more goroutines on link-heavy sites |
| One map and one lock, with "check limit and claim URL" in a single critical section | `sync.Map` plus a separate counter        | Racing workers can never overshoot `maxPages` or fetch a page twice                                                                                                                 |
| Normalised keys that keep the query string and path case                            | Stripping queries, lowercasing everything | Queries and path case can change the page, so merging them could hide real pages                                                                                                    |
| One shared HTTP client with idle connections sized to `maxConcurrency`              | A client per request                      | Connection reuse avoids repeated TLS handshakes                                                                                                                                     |
| goquery                                                                             | Raw tokenizer or regular expressions      | A real HTML5 parser copes with broken markup, and selectors keep extractors short                                                                                                   |
| Indented JSON sorted by URL                                                         | CSV or streaming output                   | Pages hold lists of links and images, which JSON represents naturally, and sorting makes uncapped runs reproducible                                                                 |
| Everything under `internal/`                                                        | Public packages                           | It is a tool, not a library, for now                                                                                                                                                |

## The website

The project site at [rajaraja.surajvijayan.dev](https://rajaraja.surajvijayan.dev) lives in [`site/`](site/). It is plain HTML, CSS and JavaScript with no build step, kept apart from the Go code. Open `site/index.html` in a browser to preview it.

It is published by GitHub Actions with [`.github/workflows/pages.yml`](.github/workflows/pages.yml) on every push to `main`. The workflow runs `go vet` and `go test` first, so a failing test stops the deploy, then publishes only the `site/` folder.

One-time custom domain setup:

1. Repository **Settings > Pages**: set **Source** to **GitHub Actions** and enter `rajaraja.surajvijayan.dev` as the custom domain.
2. At the DNS provider for `surajvijayan.dev`: add a `CNAME` record for `rajaraja` pointing to `its-me-sv.github.io`.
3. Once the certificate is issued, tick **Enforce HTTPS**.

## Contributing

Found a page it handles wrongly? [Open an issue](https://github.com/its-me-sv/rajaraja-web-crawler/issues) with:

- the command you ran,
- the URL involved (if it is public), and
- what you expected against what happened.

Pull requests are welcome. Please keep `go vet` and `go test ./...` passing and add a table test for any new extraction or normalisation rule.

## License

- Code: MIT, see [LICENSE](LICENSE).
- The Rajaraja logo and the files in `docs/brand/` are not covered by the MIT license. All rights reserved by Suraj Vijayan.

## Acknowledgements

- HTML parsing by [goquery](https://github.com/PuerkitoBio/goquery), [cascadia](https://github.com/andybalholm/cascadia) and [golang.org/x/net/html](https://pkg.go.dev/golang.org/x/net/html).
- Example crawls run against the scraping practice site at [learnwebscraping.dev](https://learnwebscraping.dev/practice/ecommerce/).
