package loadbalancer

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	network2 "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-04-01/network"
	"math/rand"
	"strings"
	"time"
)

// Flags exposes the loadbalancer reboot command arguments/flags
type Flags struct {
	Mode             string
	Scope            string
	ResourceGroup    string
	LoadBalancerName string
}

func Kill(subId string, flags Flags) error {
	c, err := NewLoadBalancersClient(subId)
	if err != nil {
		return err
	}
	lbList, err := listLBs(err, c, flags)
	if err != nil {
		return err
	}
	instanceToDelete := selectInstanceToDelete(lbList.Values(), flags)
	_, err = c.Delete(context.TODO(), flags.ResourceGroup, instanceToDelete)
	return err
}

func listLBs(err error, c network.LoadBalancersClient, flags Flags) (network2.LoadBalancerListResultPage, error) {
	lbList, err := c.List(context.TODO(), flags.ResourceGroup)
	if err != nil {
		return network2.LoadBalancerListResultPage{}, err
	}
	return lbList, nil
}

func selectInstanceToDelete(lbList []network2.LoadBalancer, flags Flags) string {
	if strings.Compare("random", flags.Mode) == 0 {
		rand.Seed(time.Now().UnixNano())
		return *lbList[rand.Intn(len(lbList))].Name
	} else {
		return *lbList[0].Name
	}
}
