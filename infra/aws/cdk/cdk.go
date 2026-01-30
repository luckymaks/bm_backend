package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
	"github.com/luckymaks/bm_backend/infra/aws"
	"github.com/luckymaks/bm_backend/infra/aws/cdk/cdkutil"
)

func main() {
	defer jsii.Close()
	app := awscdk.NewApp(nil)

	// Set to true when you have a domain configured
	enableCustomDomain := false

	// First, create shared primary region stack first
	primarySharedStack := cdkutil.NewStack(app, cdkutil.PrimaryRegion(app))
	primary := aws.NewShared(primarySharedStack, aws.SharedProps{
		EnableCustomDomain: enableCustomDomain,
	})

	// Then, create secondary shared region stacks with dependency on primary
	secondaries := map[string]aws.Shared{}
	for _, region := range cdkutil.SecondaryRegions(app) {
		secondarySharedStack := cdkutil.NewStack(app, region)
		secondaries[region] = aws.NewShared(secondarySharedStack, aws.SharedProps{
			EnableCustomDomain: enableCustomDomain,
		})
		secondarySharedStack.AddDependency(primarySharedStack, jsii.String("Primary region must deploy first"))
	}

	// Then, create stacks for the deployments.
	for _, deploymentIdent := range cdkutil.AllowedDeployments(app) {
		primaryDeploymentStack := cdkutil.NewStack(app, cdkutil.PrimaryRegion(app), deploymentIdent)

		deploymentProps := aws.DeploymentProps{
			DeploymentIdent: jsii.String(deploymentIdent),
		}
		if enableCustomDomain {
			deploymentProps.HostedZone = primary.DNS().HostedZone()
			deploymentProps.Certificate = primary.Certificate().WildcardCertificate()
		}
		aws.NewDeployment(primaryDeploymentStack, deploymentProps)
		primaryDeploymentStack.AddDependency(primarySharedStack, jsii.String("Primary shared stack must deploy first"))

		// Finally, secondary region stacks for each deployment.
		for _, region := range cdkutil.SecondaryRegions(app) {
			secondaryDeploymentStack := cdkutil.NewStack(app, region, deploymentIdent)

			secondaryDeploymentProps := aws.DeploymentProps{
				DeploymentIdent: jsii.String(deploymentIdent),
			}
			if enableCustomDomain {
				secondaryDeploymentProps.HostedZone = secondaries[region].DNS().HostedZone()
				secondaryDeploymentProps.Certificate = secondaries[region].Certificate().WildcardCertificate()
			}
			aws.NewDeployment(secondaryDeploymentStack, secondaryDeploymentProps)

			secondaryDeploymentStack.AddDependency(primaryDeploymentStack,
				jsii.String("Primary region deployment must deploy first"))
		}
	}

	app.Synth(nil)
}
