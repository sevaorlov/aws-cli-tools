package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/rds"
)

func DBInfo(sess *session.Session, rdsId *string, region *string) {
	config := &aws.Config{Region: region}
	rdsSvc := rds.New(sess, config)
	cloudwatchSvc := cloudwatch.New(sess, config)

	var describeParams *rds.DescribeDBInstancesInput
	if len(*rdsId) != 0 {
		describeParams = &rds.DescribeDBInstancesInput{DBInstanceIdentifier: rdsId}
	}

	resp, err := rdsSvc.DescribeDBInstances(describeParams)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for _, dbInstance := range resp.DBInstances {
		wg.Add(1)
		go func(cloudwatchSvc cloudwatch.CloudWatch, dbInstance rds.DBInstance) {
			defer wg.Done()
			PrintDBInstanceInfo(&cloudwatchSvc, &dbInstance)
		}(*cloudwatchSvc, *dbInstance)
	}
	wg.Wait()
}

func GetMetricStatisticsInputForDBInstance(instanceId *string, metricName string, statistics *string, period int64, duration time.Duration) *cloudwatch.GetMetricStatisticsInput {
	return &cloudwatch.GetMetricStatisticsInput{
		EndTime:    aws.Time(time.Now()),                    // Required
		MetricName: aws.String(metricName),                  // Required
		Namespace:  aws.String("AWS/RDS"),                   // Required
		Period:     aws.Int64(period),                       // Required
		StartTime:  aws.Time(time.Now().Add(-1 * duration)), // Required
		Statistics: []*string{ // Required
			statistics, // Required
		},
		Dimensions: []*cloudwatch.Dimension{
			{ // Required
				Name:  aws.String("DBInstanceIdentifier"), // Required
				Value: instanceId,                         // Required
			},
			// More values...

		},
	}
}

func GetDBInstanceMetrics(svc *cloudwatch.CloudWatch, metricName string, instanceId *string, statistics *string, period int64, duration time.Duration) *cloudwatch.GetMetricStatisticsOutput {
	params := GetMetricStatisticsInputForDBInstance(instanceId, metricName, statistics, period, duration)
	resp, err := svc.GetMetricStatistics(params)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return resp
}

func PrintDBInstanceInfo(cloudwatchSvc *cloudwatch.CloudWatch, dbInstance *rds.DBInstance) {
	statistics := "Average"

	cpuData := GetDBInstanceMetrics(cloudwatchSvc, "CPUUtilization", dbInstance.DBInstanceIdentifier, &statistics, 60, 10*time.Minute)
	writeIOPSData := GetDBInstanceMetrics(cloudwatchSvc, "WriteIOPS", dbInstance.DBInstanceIdentifier, &statistics, 60, 10*time.Minute)
	readIOPSData := GetDBInstanceMetrics(cloudwatchSvc, "ReadIOPS", dbInstance.DBInstanceIdentifier, &statistics, 60, 10*time.Minute)
	memoryData := GetDBInstanceMetrics(cloudwatchSvc, "FreeableMemory", dbInstance.DBInstanceIdentifier, &statistics, 60, time.Minute)
	storageSpaceData := GetDBInstanceMetrics(cloudwatchSvc, "FreeStorageSpace", dbInstance.DBInstanceIdentifier, &statistics, 60, time.Minute)

	output := []string{}
	output = append(output, *dbInstance.DBInstanceIdentifier)
	output = append(output, *dbInstance.DBInstanceStatus)
	output = append(output, *dbInstance.DBInstanceClass)

	output = append(output, *FormatData(cpuData, "", &statistics))
	output = append(output, *FormatData(writeIOPSData, "", &statistics))
	output = append(output, *FormatData(readIOPSData, "", &statistics))
	output = append(output, *FormatData(memoryData, "GB", &statistics))
	output = append(output, *FormatData(storageSpaceData, "GB", &statistics))

	output = append(output, "---")

	fmt.Println(strings.Join(output, "\n"))
}

func FormatData(data *cloudwatch.GetMetricStatisticsOutput, format string, statistics *string) *string {
	str := *data.Label + ": "
	for _, datapoint := range data.Datapoints {
		var value float64

		if *statistics == "Average" {
			value = *datapoint.Average
		}

		if format == "GB" {
			value = *datapoint.Average / 1024.0 / 1024.0 / 1024.0
		}

		str += strconv.FormatFloat(value, 'f', 1, 64) + format + " "
	}

	return &str
}
