package main

import (
	"flag"

	"github.com/aws/aws-sdk-go/aws/session"
)

func main() {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	command := flag.String("command", "", "command to execute")

	// dbinfo
	rdsId := flag.String("rds", "", "rds identifier")
	region := flag.String("reg", "eu-west-1", "aws region")

	// ssh
	stackNamePrefix := flag.String("prefix", "", "stack name prefix")
	stackName := flag.String("stack", "", "stack name (without prefix optionally)")
	layerShortname := flag.String("layer", "", "layer short name")

	flag.Parse()

	if *command == "ssh" {
		InstancesSSH(sess, stackName, layerShortname, stackNamePrefix)
	}

	if *command == "dbinfo" {
		DBInfo(sess, rdsId, region)
	}
}
