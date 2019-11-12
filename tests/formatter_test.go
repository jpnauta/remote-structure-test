package tests

import (
	"bytes"
	"github.com/jpnauta/remote-structure-test/pkg/color"
	"io"
	"testing"
)

func compareText(t *testing.T, expected, actual string, expectedN int, actualN int, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("did not expect error when formatting text but got %s", err)
	}
	if actualN != expectedN {
		t.Errorf("expected formatter to have written %d bytes but wrote %d", expectedN, actualN)
	}
	if actual != expected {
		t.Errorf("formatting not applied to text. Expected \"%s\" but got \"%s\"", expected, actual)
	}
}

func TestFprint(t *testing.T) {
	defer func(f func(io.Writer) bool) { color.ColoredOutput = f }(color.ColoredOutput)
	color.ColoredOutput = func(io.Writer) bool { return true }

	var b bytes.Buffer
	n, err := color.Green.Fprint(&b, "It's not easy being")
	expected := "\033[32mIt's not easy being\033[0m"
	compareText(t, expected, b.String(), 28, n, err)
}

func TestFprintln(t *testing.T) {
	defer func(f func(io.Writer) bool) { color.ColoredOutput = f }(color.ColoredOutput)
	color.ColoredOutput = func(io.Writer) bool { return true }

	var b bytes.Buffer
	n, err := color.Green.Fprintln(&b, "2", "less", "chars!")
	expected := "\033[32m2 less chars!\033[0m\n"
	compareText(t, expected, b.String(), 23, n, err)
}

func TestFprintf(t *testing.T) {
	defer func(f func(io.Writer) bool) { color.ColoredOutput = f }(color.ColoredOutput)
	color.ColoredOutput = func(io.Writer) bool { return true }

	var b bytes.Buffer
	n, err := color.Green.Fprintf(&b, "It's been %d %s", 1, "week")
	expected := "\033[32mIt's been 1 week\033[0m"
	compareText(t, expected, b.String(), 25, n, err)
}

func TestFprintNoTTY(t *testing.T) {
	var b bytes.Buffer
	expected := "It's not easy being"
	n, err := color.Green.Fprint(&b, expected)
	compareText(t, expected, b.String(), 19, n, err)
}

func TestFprintlnNoTTY(t *testing.T) {
	var b bytes.Buffer
	n, err := color.Green.Fprintln(&b, "2", "less", "chars!")
	expected := "2 less chars!\n"
	compareText(t, expected, b.String(), 14, n, err)
}

func TestFprintfNoTTY(t *testing.T) {
	var b bytes.Buffer
	n, err := color.Green.Fprintf(&b, "It's been %d %s", 1, "week")
	expected := "It's been 1 week"
	compareText(t, expected, b.String(), 16, n, err)
}

func TestFprintTTYNoColor(t *testing.T) {
	defer func(f func(io.Writer) bool) { color.IsTerminal = f }(color.IsTerminal)
	color.IsTerminal = func(io.Writer) bool { return true }
	defer func() { color.NoColor = false }()
	color.NoColor = true

	var b bytes.Buffer
	expected := "It's not easy being"
	n, err := color.Green.Fprint(&b, expected)
	compareText(t, expected, b.String(), 19, n, err)
}

func TestFprintlnTTYNoColor(t *testing.T) {
	defer func(f func(io.Writer) bool) { color.IsTerminal = f }(color.IsTerminal)
	color.IsTerminal = func(io.Writer) bool { return true }
	defer func() { color.NoColor = false }()
	color.NoColor = true

	var b bytes.Buffer
	n, err := color.Green.Fprintln(&b, "2", "less", "chars!")
	expected := "2 less chars!\n"
	compareText(t, expected, b.String(), 14, n, err)
}

func TestFprintfTTYNoColor(t *testing.T) {
	defer func(f func(io.Writer) bool) { color.IsTerminal = f }(color.IsTerminal)
	color.IsTerminal = func(io.Writer) bool { return true }
	defer func() { color.NoColor = false }()
	color.NoColor = true

	var b bytes.Buffer
	n, err := color.Green.Fprintf(&b, "It's been %d %s", 1, "week")
	expected := "It's been 1 week"
	compareText(t, expected, b.String(), 16, n, err)
}

func TestOverwriteDefault(t *testing.T) {
	CheckDeepEqual(t, color.None, color.Default)
	color.OverwriteDefault(color.Red)
	CheckDeepEqual(t, color.Red, color.Default)
}
