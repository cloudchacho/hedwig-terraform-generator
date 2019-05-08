package main

import (
	"errors"
	"gopkg.in/urfave/cli.v1"
)

func newConfig(c *cli.Context, filename string) (interface{}, error) {
	if c.GlobalString(cloudProviderFlag) == cloudProviderGoogle {
		return newGoogleConfig(filename)
	} else if c.GlobalString(cloudProviderFlag) == cloudProviderAWS {
		return newAWSConfig(filename)
	} else {
		return nil, errors.New("invalid cloud provider")
	}
}
