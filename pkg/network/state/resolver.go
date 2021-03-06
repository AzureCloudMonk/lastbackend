//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package state

import (
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"sync"
)

type ResolverState struct {
	lock      sync.RWMutex
	resolvers map[string]*types.ResolverManifest
}

func (n *ResolverState) GetResolvers() map[string]*types.ResolverManifest {
	return n.resolvers
}

func (n *ResolverState) AddResolver(cidr string, sn *types.ResolverManifest) {
	log.V(logLevel).Debugf("Stage: ResolverState: add resolver: %s", cidr)
	n.SetResolver(cidr, sn)
}

func (n *ResolverState) SetResolver(cidr string, sn *types.ResolverManifest) {
	log.V(logLevel).Debugf("Stage: ResolverState: set resolver: %s", cidr)
	n.lock.Lock()
	defer n.lock.Unlock()

	if _, ok := n.resolvers[cidr]; ok {
		delete(n.resolvers, cidr)
	}

	n.resolvers[cidr] = sn
}

func (n *ResolverState) GetResolver(cidr string) *types.ResolverManifest {
	log.V(logLevel).Debugf("Stage: ResolverState: get resolver: %s", cidr)
	n.lock.Lock()
	defer n.lock.Unlock()
	s, ok := n.resolvers[cidr]
	if !ok {
		return nil
	}
	return s
}

func (n *ResolverState) DelResolver(cidr string) {
	log.V(logLevel).Debugf("Stage: ResolverState: del resolver: %v", cidr)
	n.lock.Lock()
	defer n.lock.Unlock()
	if _, ok := n.resolvers[cidr]; ok {
		delete(n.resolvers, cidr)
	}
}
