package main

import (
	"github.com/sirupsen/logrus"

	"github.com/jpnauta/remote-structure-test/cmd/remote-structure-test/app"
)

func main() {
	if err := app.Run(); err != nil {
		logrus.Fatal(err)
	}
}
