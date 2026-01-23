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
	
	// First, create shared primary region stack first
	primarySharedStack := cdkutil.NewStack(app, cdkutil.PrimaryRegion(app))
	primary := aws.NewShared(primarySharedStack, aws.SharedProps{})
	
	// Then, create secondary shared region stacks with dependency on primary
	secondaries := map[string]aws.Shared{}
	for _, region := range cdkutil.SecondaryRegions(app) {
		secondarySharedStack := cdkutil.NewStack(app, region)
		secondaries[region] = aws.NewShared(secondarySharedStack, aws.SharedProps{})
		secondarySharedStack.AddDependency(primarySharedStack, jsii.String("Primary region must deploy first"))
	}
	
	// Then, create stacks for the deployments.
	for _, deploymentIdent := range cdkutil.AllowedDeployments(app) {
		primaryDeploymentStack := cdkutil.NewStack(app, cdkutil.PrimaryRegion(app), deploymentIdent)
		aws.NewDeployment(primaryDeploymentStack, aws.DeploymentProps{
			DeploymentIdent: jsii.String(deploymentIdent),
			HostedZone:      primary.DNS().HostedZone(),
			Certificate:     primary.Certificate().WildcardCertificate(),
			Identity:        primary.Identity(),
			CrewIdentity:    primary.CrewIdentity(),
		})
		primaryDeploymentStack.AddDependency(primarySharedStack, jsii.String("Primary shared stack must deploy first"))
		
		// Finally, secondary region stacks for each deployment.
		for _, region := range cdkutil.SecondaryRegions(app) {
			secondaryDeploymentStack := cdkutil.NewStack(app, region, deploymentIdent)
			aws.NewDeployment(secondaryDeploymentStack, aws.DeploymentProps{
				DeploymentIdent: jsii.String(deploymentIdent),
				HostedZone:      secondaries[region].DNS().HostedZone(),
				Certificate:     secondaries[region].Certificate().WildcardCertificate(),
				Identity:        secondaries[region].Identity(),
				CrewIdentity:    secondaries[region].CrewIdentity(),
			})
			
			secondaryDeploymentStack.AddDependency(primaryDeploymentStack,
				jsii.String("Primary region deployment must deploy first"))
		}
	}
	
	app.Synth(nil)
}
