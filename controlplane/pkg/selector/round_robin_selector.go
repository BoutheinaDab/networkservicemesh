package selector

import (
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/local/connection"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/registry"
)

type roundRobinSelector struct {
	sync.Mutex
	roundRobin map[string]int
}

func NewRoundRobinSelector() Selector {
	return &roundRobinSelector{
		roundRobin: make(map[string]int),
	}
}

func (rr *roundRobinSelector) SelectEndpoint(requestConnection *connection.Connection, ns *registry.NetworkService, networkServiceEndpoints []*registry.NetworkServiceEndpoint) *registry.NetworkServiceEndpoint {
        
        logrus.Infof("start RoundRobin for ns.name %s and requestConnection.getid %s ", ns.GetName(), requestConnection.GetId())
	if rr == nil {
		return nil
	}
	if len(networkServiceEndpoints) == 0 {
		return nil
	}
	rr.Lock()
	defer rr.Unlock()
	idx := rr.roundRobin[ns.GetName()] % len(networkServiceEndpoints)
	endpoint := networkServiceEndpoints[idx]
        logrus.Infof("RoundRobin selected %v idx %d ", endpoint, idx)
	if endpoint == nil {
		return nil
	}
	rr.roundRobin[ns.GetName()] = rr.roundRobin[ns.GetName()] + 1
	return endpoint
}
