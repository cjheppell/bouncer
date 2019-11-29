package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type awsFirewallClient struct {
	sess            *session.Session
	securityGroupID string
}

// NewFirewallClient creates a new AWS Firewall Client for manipulating AWS security group rules
func NewFirewallClient(region, securityGroupID string) (awsFirewallClient, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	return awsFirewallClient{
		sess: sess,
		securityGroupID: securityGroupID,
	}, nil
}

// ExposePort exposes the specified port to the public in an AWS security group
func (pe awsFirewallClient) ExposePort(port int32) error {
	ec2Svc := ec2.New(pe.sess)

	newIngressRule := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: aws.String(pe.securityGroupID),
		IpPermissions: []*ec2.IpPermission{
			{
				FromPort:   aws.Int64(int64(port)),
				IpProtocol: aws.String("tcp"),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp:      aws.String("0.0.0.0/0"),
						Description: aws.String(fmt.Sprintf("Public access to NodePort %d", port)),
					},
				},
				ToPort: aws.Int64(int64(port)),
			},
		},
	}

	if _, err := ec2Svc.AuthorizeSecurityGroupIngress(newIngressRule); err != nil {
		return err
	}

	return nil
}

// UnexposePort completely unexposes the specified port in an AWS security group
func (pe awsFirewallClient) UnexposePort(port int32) error {
	ec2Svc := ec2.New(pe.sess)

	ingressRule := &ec2.RevokeSecurityGroupIngressInput{
		GroupId: aws.String(pe.securityGroupID),
		IpPermissions: []*ec2.IpPermission{
			{
				FromPort:   aws.Int64(int64(port)),
				IpProtocol: aws.String("tcp"),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp:      aws.String("0.0.0.0/0"),
						Description: aws.String(fmt.Sprintf("Public access to NodePort %d", port)),
					},
				},
				ToPort: aws.Int64(int64(port)),
			},
		},
	}

	if _, err := ec2Svc.RevokeSecurityGroupIngress(ingressRule); err != nil {
		return err
	}

	return nil
}