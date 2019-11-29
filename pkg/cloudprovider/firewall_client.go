package cloudprovider

import (
	"fmt"
	"github.com/cjheppell/bouncer/pkg/cloudprovider/aws"
	"github.com/cjheppell/bouncer/pkg/cloudprovider/noop"
	"os"
)

// FirewallClient exposes an interface for maniuplating firewall rules for an appropriate cloud provider
type FirewallClient interface {
	ExposePort(port int32) error
	UnexposePort(port int32) error
}

// GetFirewallClient gets an appropriate firewall client based on the environment variables running this process
func GetFirewallClient() (FirewallClient, error){
	region := os.Getenv("AWS_REGION")
	if region != "" {
		return getAwsClient()
	}

	return getNoOpClient()
}


func getAwsClient() (FirewallClient, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		return nil, fmt.Errorf("AWS_REGION env var was not set")
	}

	securityGroupID := os.Getenv("AWS_SECURITY_GROUP_ID")
	if securityGroupID == "" {
		return nil, fmt.Errorf("AWS_SECURITY_GROUP_ID env var was not set")
	}

	firewallClient, err := aws.NewFirewallClient(region, securityGroupID)
	if err != nil {
		return nil, err
	}

	return firewallClient, nil
}

func getNoOpClient() (FirewallClient, error) {
	return noop.NoopFirewallClient{}, nil
}
