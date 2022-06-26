package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	apigwv2 "github.com/aws/aws-cdk-go/awscdkapigatewayv2alpha/v2"
	integrations "github.com/aws/aws-cdk-go/awscdkapigatewayv2integrationsalpha/v2"
	lambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkStackProps struct {
	awscdk.StackProps
}

func NewCdkStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	table := awsdynamodb.NewTable(stack, jsii.String("Table"), &awsdynamodb.TableProps{
		BillingMode:   awsdynamodb.BillingMode_PROVISIONED,
		ReadCapacity:  jsii.Number(1),
		WriteCapacity: jsii.Number(1),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("PK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: jsii.String("SK"),
			Type: awsdynamodb.AttributeType_STRING,
		},
	})

	fn := lambdago.NewGoFunction(stack, jsii.String("ApiLambda"), &lambdago.GoFunctionProps{
		Entry: jsii.String("lambda"),
		Environment: &map[string]*string{
			"FARMFA_DYNAMODB_TABLE": table.TableName(),
		},
		Architecture: awslambda.Architecture_ARM_64(),
		LogRetention: awslogs.RetentionDays_ONE_WEEK,
	})
	table.GrantReadWriteData(fn)

	apigwv2.NewHttpApi(stack, jsii.String("HttpApiGateway"), &apigwv2.HttpApiProps{
		DefaultIntegration: integrations.NewHttpLambdaIntegration(jsii.String("LambdaIntegration"), fn, nil),
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewCdkStack(app, "FarmfaStack", &CdkStackProps{
		awscdk.StackProps{},
	})

	app.Synth(nil)
}
