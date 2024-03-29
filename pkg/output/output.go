package output

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	color "github.com/jpnauta/remote-structure-test/pkg/color"
	types "github.com/jpnauta/remote-structure-test/pkg/types/unversioned"
)

var bannerLength = 27 // default banner length

func OutputResult(out io.Writer, result *types.TestResult) {
	color.Default.Fprintf(out, "=== RUN: %s\n", result.Name)
	if result.Pass {
		color.Green.Fprintln(out, "--- PASS")
	} else {
		color.Red.Fprintln(out, "--- FAIL")
	}
	color.Default.Fprintf(out, "duration: %s\n", result.Duration.String())
	if result.Stdout != "" {
		color.Blue.Fprintf(out, "stdout: %s\n", result.Stdout)
	}
	if result.Stderr != "" {
		color.Blue.Fprintf(out, "stderr: %s\n", result.Stderr)
	}
	for _, s := range result.Errors {
		color.Yellow.Fprintf(out, "Error: %s\n", s)
	}
}

func Banner(out io.Writer, filename string) {
	fileStr := fmt.Sprintf("====== Test file: %s ======", filepath.Base(filename))
	bannerLength = len(fileStr)
	color.Purple.Fprintln(out, "\n"+strings.Repeat("=", bannerLength))
	color.Purple.Fprintln(out, fileStr)
	color.Purple.Fprintln(out, strings.Repeat("=", bannerLength))
}

func FinalResults(out io.Writer, jsonOut bool, result types.SummaryObject) error {
	if jsonOut {
		res, err := json.Marshal(result)
		if err != nil {
			return errors.Wrap(err, "marshalling json")
		}
		res = append(res, []byte("\n")...)
		_, err = out.Write(res)
		return err
	}
	if bannerLength%2 == 0 {
		bannerLength++
	}
	if result.Total == 0 {
		color.Red.Fprintln(out, "No tests run! Check config file format.")
		return nil
	}
	color.Default.Fprintln(out, "")
	color.Default.Fprintln(out, strings.Repeat("=", bannerLength))
	color.Default.Fprintf(out, "%s RESULTS %s\n", strings.Repeat("=", (bannerLength-9)/2), strings.Repeat("=", (bannerLength-9)/2))
	color.Default.Fprintln(out, strings.Repeat("=", bannerLength))
	color.LightGreen.Fprintf(out, "Passes:      %d\n", result.Pass)
	color.LightRed.Fprintf(out, "Failures:    %d\n", result.Fail)
	color.Default.Fprintf(out, "Duration:    %s\n", result.Duration.String())
	color.Cyan.Fprintf(out, "Total tests: %d\n", result.Total)
	color.Default.Fprintln(out, "")
	if result.Fail == 0 {
		color.Green.Fprintln(out, "PASS")
	} else {
		color.Red.Fprintln(out, "FAIL")
	}
	return nil
}
