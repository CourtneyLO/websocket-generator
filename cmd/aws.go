package cmd

import (
		"github.com/aws/aws-sdk-go/aws"
		"github.com/aws/aws-sdk-go/aws/session"
		"github.com/aws/aws-sdk-go/service/ssm"
)

func GetAWSParameter(parameterName string, awsRegion string) (string, error) {
	sess := session.Must(session.NewSession(&aws.Config{
    Region: aws.String(awsRegion),
	}))
	ssmsvc := ssm.New(sess)
	param, error := ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name: aws.String(parameterName),
	})

	value := *param.Parameter.Value

	return value, error
}
