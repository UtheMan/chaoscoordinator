package loadbalancer

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure"
)

const userAgent string = "genesys"

func NewLoadBalancersClient(subId string) (network.LoadBalancersClient, error) {
	client, err := newLoadBalancersClient(subId)
	if err != nil {
		return network.LoadBalancersClient{}, err
	}
	return client, nil
}

func newLoadBalancersClient(subID string) (network.LoadBalancersClient, error) {
	a, err := azure.InjectAuthorizer()
	if err != nil {
		return network.LoadBalancersClient{}, err
	}
	client := network.NewLoadBalancersClient(subID)
	client.Authorizer = a
	if err := client.AddToUserAgent(userAgent); err != nil {
		return network.LoadBalancersClient{}, err
	}
	return client, nil
}
