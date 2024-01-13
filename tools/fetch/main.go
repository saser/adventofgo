// Binary fetch sends HTTP requests to https://adventofcode.com to fetch problem
// inputs and any available existing answers. The results are written out to
// files in a given output directory.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/golang/glog"
	"golang.org/x/net/publicsuffix"
)

var (
	session   = flag.String("session", "", `The value of the "session" cookie needed to authenticate to https://adventofcode.com. Grab it from your browser's cookie store.`)
	year      = flag.Int("year", 2015, "The event's year.")
	day       = flag.Int("day", 1, "The event's day.")
	outputDir = flag.String("output_dir", "", "Path to a directory in which to write output files. File names will have the form <output_dir>/year<year>_day<day>_*.")
)

func aocURL() *url.URL {
	const base = "https://adventofcode.com"
	u, err := url.Parse(base)
	if err != nil {
		panic(fmt.Errorf("parse %q as URL: %v", base, err))
	}
	return u
}

var _ = aocURL()

// buildHTTPClient creates a *http.Client with a cookie set for
// https://adventofcode.com that sets the "session" key to the given value.
func buildHTTPClient(session string) (*http.Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, fmt.Errorf("create cookie jar: %w", err)
	}
	glog.V(1).Infof(`Setting "session" cookie to %q.`, session)
	jar.SetCookies(aocURL(), []*http.Cookie{{Name: "session", Value: session}})
	return &http.Client{Jar: jar}, nil
}

// getInput issues a HTTP GET request for https://adventofcode/<year>/day/<day>/input.
func getInput(ctx context.Context, c *http.Client, year int, day int) (string, error) {
	u := aocURL()
	u.Path = path.Join(fmt.Sprint(year), "day", fmt.Sprint(day), "input")
	glog.V(1).Infof("Fetching input from %q.", u.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("GET input page for year %d, day %d: build HTTP GET request: %v", year, day, err)
	}
	res, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("GET input page for year %d, day %d: do HTTP GET request: %v", year, day, err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("GET input page for year %d, day %d: read response body: %v", year, day, err)
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GET input page for year %d, day %d: HTTP GET request returned status %q and response body: %s", year, day, res.Status, string(body))
	}
	glog.V(1).Infof("Fetched %d bytes of input.", len(body))
	return string(body), nil
}

var answerRE = regexp.MustCompile(`Your puzzle answer was <code>(.+?)</code>`)

// getAnswers issues a HTTP GET request for
// https://adventofcode/<year>/day/<day>. If the returned body contains answers
// (because the problem has been solved before) those are parsed out using a
// rudimentary regex. If neither part 1 or part 2 has been solved, this function
// returns two empty strings with a nil error.
func getAnswers(ctx context.Context, c *http.Client, year int, day int) (part1 string, part2 string, err error) {
	u := aocURL()
	u.Path = path.Join(fmt.Sprint(year), "day", fmt.Sprint(day))
	glog.V(1).Infof("Parsing answers for year %d, day %d from %q.", year, day, u.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", "", fmt.Errorf("GET problem page for year %d, day %d: build HTTP GET request: %v", year, day, err)
	}
	res, err := c.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("GET problem page for year %d, day %d: do HTTP GET request: %v", year, day, err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", fmt.Errorf("GET problem page for year %d, day %d: read response body: %v", year, day, err)
	}
	if res.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("GET problem page for year %d, day %d: HTTP GET request returned status %q and response body: %s", year, day, res.Status, string(body))
	}
	strBody := string(body)
	if strings.Contains(strBody, "To play, please identify yourself via one of these services") {
		return "", "", fmt.Errorf("GET problem page for year %d, day %d: problem page is for users not logged in; is the session cookie correct?", year, day)
	}
	matches := answerRE.FindAllStringSubmatch(string(body), 2)
	if matches == nil {
		glog.V(1).Infof("Found no answers for year %d, day %d.", year, day)
		return "", "", nil
	}
	if len(matches) >= 1 {
		part1 = matches[0][1]
	}
	if len(matches) >= 2 {
		part2 = matches[1][1]
	}
	return part1, part2, nil
}

// dataset represents all the data we could gather from the website for a given
// year and day.
type dataset struct {
	Year, Day    int
	Input        string
	Part1, Part2 string
}

// writeDataset writes out the information in a dataset to files in the given
// directory. It ensures that all written files have a trailing newline.
func writeDataset(ds dataset, dir string) error {
	base := filepath.Join(dir, fmt.Sprintf("year%d_day%02d", ds.Year, ds.Day))
	glog.V(1).Infof("Using %q as the base for filenames.", base)

	ensureNewline := func(s string) string { return strings.TrimSpace(s) + "\n" }

	inputPath := base + "_input"
	glog.V(1).Infof("Writing input file to %q.", inputPath)
	if err := os.WriteFile(inputPath, []byte(ensureNewline(ds.Input)), fs.FileMode(0644)); err != nil {
		return fmt.Errorf("write input file: %v", err)
	}

	part1Path := base + "_part1_output"
	if ds.Part1 != "" {
		glog.V(1).Infof("Writing part 1 answer (%q) to %q.", ds.Part1, part1Path)
	} else {
		glog.V(1).Infof("Writing empty part 1 answer to %q.", part1Path)
	}
	if err := os.WriteFile(part1Path, []byte(ensureNewline(ds.Part1)), fs.FileMode(0644)); err != nil {
		return fmt.Errorf("write part 1 answer file: %v", err)
	}

	part2Path := base + "_part2_output"
	if ds.Part2 != "" {
		glog.V(1).Infof("Writing part 2 answer (%q) to %q.", ds.Part2, part2Path)
	} else {
		glog.V(1).Infof("Writing empty part 2 answer to %q.", part2Path)
	}
	glog.V(1).Infof("Writing part 2 file to %q.", part2Path)
	if err := os.WriteFile(part2Path, []byte(ensureNewline(ds.Part2)), fs.FileMode(0644)); err != nil {
		return fmt.Errorf("write part 2 answer file: %v", err)
	}

	return nil
}

func errmain() error {
	ctx := context.Background()

	if *session == "" {
		return errors.New("-session is required")
	}
	if *year < 2015 {
		return fmt.Errorf("-year=%d is invalid; needs to be 2015 or higher", *year)
	}
	if *day < 1 || *day > 25 {
		return fmt.Errorf("-day=%d is invalid; needs to be in [1, 25]", *day)
	}
	if *outputDir == "" {
		return errors.New("-output_dir is required")
	}

	c, err := buildHTTPClient(*session)
	if err != nil {
		return fmt.Errorf("build HTTP client: %v", err)
	}

	input, err := getInput(ctx, c, *year, *day)
	if err != nil {
		return fmt.Errorf("fetch input: %v", err)
	}

	part1, part2, err := getAnswers(ctx, c, *year, *day)
	if err != nil {
		return fmt.Errorf("fetch answers: %v", err)
	}
	if part1 != "" {
		glog.V(1).Infof("Answer to part 1 is %q.", part1)
	}
	if part2 != "" {
		glog.V(1).Infof("Answer to part 2 is %q.", part2)
	}

	ds := dataset{
		Year:  *year,
		Day:   *day,
		Input: input,
		Part1: part1,
		Part2: part2,
	}
	if err := writeDataset(ds, *outputDir); err != nil {
		return fmt.Errorf("write dataset: %v", err)
	}

	glog.Infof("Wrote files for year %d, day %d.", *year, *day)

	return nil
}

func main() {
	flag.Parse()
	if err := errmain(); err != nil {
		glog.Exit(err)
	}
}
