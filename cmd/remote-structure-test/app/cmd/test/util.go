package test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/jpnauta/remote-structure-test/pkg/config"
	"github.com/jpnauta/remote-structure-test/pkg/drivers"
	"github.com/jpnauta/remote-structure-test/pkg/output"
	"github.com/jpnauta/remote-structure-test/pkg/types"
	"github.com/jpnauta/remote-structure-test/pkg/types/unversioned"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func ValidateArgs(opts *config.StructureTestOptions) error {
	if opts.Host == "" {
		return fmt.Errorf("Please supply host to run tests on")
	}
	if len(opts.ConfigFiles) == 0 {
		return fmt.Errorf("Please provide at least one test config file")
	}
	return nil
}

func Parse(fp string, args *drivers.DriverConfig, driverImpl func(drivers.DriverConfig) (drivers.Driver, error)) (types.StructureTest, error) {
	testContents, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	// We first have to unmarshal to determine the schema version, then we unmarshal again
	// to do the full parse.
	var unmarshal types.Unmarshaller
	var strictUnmarshal types.Unmarshaller
	var versionHolder types.SchemaVersion

	switch {
	case strings.HasSuffix(fp, ".json"):
		unmarshal = json.Unmarshal
		strictUnmarshal = json.Unmarshal
	case strings.HasSuffix(fp, ".yaml"):
		unmarshal = yaml.Unmarshal
		strictUnmarshal = yaml.UnmarshalStrict
	case strings.HasSuffix(fp, ".yml"):
		unmarshal = yaml.Unmarshal
		strictUnmarshal = yaml.UnmarshalStrict
	default:
		return nil, errors.New("Please provide valid JSON or YAML config file")
	}

	if err := unmarshal(testContents, &versionHolder); err != nil {
		return nil, err
	}

	version := versionHolder.SchemaVersion
	if version == "" {
		return nil, errors.New("Please provide JSON schema version")
	}

	var st types.StructureTest
	if schemaVersion, ok := types.SchemaVersions[version]; ok {
		st = schemaVersion()
	} else {
		return nil, errors.New("Unsupported schema version: " + version)
	}

	if err = strictUnmarshal(testContents, st); err != nil {
		return nil, errors.New("error unmarshalling config: " + err.Error())
	}

	tests, _ := st.(types.StructureTest) //type assertion
	tests.SetDriverImpl(driverImpl, *args)
	return tests, nil
}

func ProcessResults(out io.Writer, json bool, c chan interface{}) error {
	totalPass := 0
	totalFail := 0
	totalDuration := time.Duration(0)
	errStrings := make([]string, 0)
	results, err := channelToSlice(c)
	if err != nil {
		return errors.Wrap(err, "reading results from channel")
	}
	for _, r := range results {
		if !json {
			// output individual results if we're not in json mode
			output.OutputResult(out, r)
		}
		if r.IsPass() {
			totalPass++
		} else {
			totalFail++
		}
		totalDuration += r.Duration
	}
	if totalPass+totalFail == 0 || totalFail > 0 {
		errStrings = append(errStrings, "FAIL")
	}
	if len(errStrings) > 0 {
		err = fmt.Errorf(strings.Join(errStrings, "\n"))
	}

	summary := unversioned.SummaryObject{
		Total:    totalFail + totalPass,
		Pass:     totalPass,
		Fail:     totalFail,
		Duration: totalDuration,
	}
	if json {
		// only output results here if we're in json mode
		summary.Results = results
	}
	output.FinalResults(out, json, summary)

	return err
}

func channelToSlice(c chan interface{}) ([]*unversioned.TestResult, error) {
	results := []*unversioned.TestResult{}
	for elem := range c {
		elem, ok := elem.(*unversioned.TestResult)
		if !ok {
			return nil, fmt.Errorf("unexpected value found in channel: %v", elem)
		}
		results = append(results, elem)
	}
	return results, nil
}
