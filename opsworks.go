package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/opsworks"
)

func InstancesSSH(sess *session.Session, stackName *string, layerShortname *string, stackNamePrefix *string) {
	svc := opsworks.New(sess, &aws.Config{Region: aws.String("us-east-1")})

	stacks := FetchStacks(svc)

	var stack *opsworks.Stack
	var layer *opsworks.Layer
	var instance *opsworks.Instance

	if len(*stackName) != 0 {
		stack = FindStackByName(stacks, stackName, stackNamePrefix)
		if stack == nil {
			fmt.Println("Stack with provided name was not found.")
		}
	}

	if stack == nil {
		stack = ChooseStackDialog(stacks, stackNamePrefix)
	}

	layers := FetchLayersByStack(svc, stack)

	if len(*layerShortname) != 0 {
		layer = FindLayerByName(layers, *layerShortname)
		if layer == nil {
			fmt.Println("Layer with provided name was not found.")
		}
	}

	if layer == nil {
		layer = ChooseLayerDialog(layers)
	}

	instances := FetchInstancesByLayer(svc, layer)
	if len(instances) > 1 {
		instance = ChooseInstanceDialog(instances)
	} else {
		instance = instances[0]
	}

	cmd := exec.Command("ssh", "ubuntu@"+*instance.PublicIp)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func FindStackByName(stacks []*opsworks.Stack, name *string, prefix *string) *opsworks.Stack {
	for _, stack := range stacks {
		if *name == ShortStackName(stack.Name, prefix) {
			return stack
		}
	}

	return nil
}

func ShortStackName(name *string, prefix *string) string {
	res := strings.ToLower(*name)

	if len(*prefix) != 0 {
		res = strings.TrimPrefix(res, *prefix)
	}

	res = strings.Trim(res, " ")
	res = strings.Replace(res, " ", "_", -1)

	return res
}

func FetchStacks(svc *opsworks.OpsWorks) []*opsworks.Stack {
	resp, err := svc.DescribeStacks(nil)

	if err != nil {
		fmt.Println(err.Error())
		return []*opsworks.Stack{}
	}

	return resp.Stacks
}

func ChooseStackDialog(stacks []*opsworks.Stack, prefix *string) *opsworks.Stack {
	fmt.Println("Choose a stack:")

	for index, stack := range stacks {
		fmt.Println(strconv.Itoa(index+1) + ") " + ShortStackName(stack.Name, prefix))
	}

	var index int
	_, err := fmt.Scanf("%d", &index)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return stacks[index-1]
}

func FetchLayersByStack(svc *opsworks.OpsWorks, stack *opsworks.Stack) []*opsworks.Layer {
	layersInput := &opsworks.DescribeLayersInput{StackId: stack.StackId}
	resp, err := svc.DescribeLayers(layersInput)

	if err != nil {
		panic(err)
	}

	return resp.Layers
}

func FindLayerByName(layers []*opsworks.Layer, layerShortname string) *opsworks.Layer {
	for _, layer := range layers {
		if layerShortname == *layer.Shortname {
			return layer
		}
	}

	return nil
}

func ChooseLayerDialog(layers []*opsworks.Layer) *opsworks.Layer {
	fmt.Println("Choose a layer:")

	for index, layer := range layers {
		fmt.Println(strconv.Itoa(index+1) + ") " + *layer.Shortname)
	}

	var index int
	_, err := fmt.Scanf("%d", &index)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return layers[index-1]
}

func FetchInstancesByLayer(svc *opsworks.OpsWorks, layer *opsworks.Layer) []*opsworks.Instance {
	params := &opsworks.DescribeInstancesInput{LayerId: layer.LayerId}
	resp, err := svc.DescribeInstances(params)

	if err != nil {
		panic(err)
	}

	return resp.Instances
}

func ChooseInstanceDialog(instances []*opsworks.Instance) *opsworks.Instance {
	fmt.Println("Choose an instance:")

	for index, instance := range instances {
		fmt.Println(strconv.Itoa(index+1) + ") " + *instance.Hostname)
	}

	var index int
	_, err := fmt.Scanf("%d", &index)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return instances[index-1]
}
