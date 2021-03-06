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

package cache

import (
	"sync"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const logCacheIngress = "api:cache:ingress"

type CacheIngressManifest struct {
	lock      sync.RWMutex
	ingress   map[string]*types.Ingress
	discovery map[string]*types.Discovery
	routes    map[string]*types.RouteManifest
	manifests map[string]*types.IngressManifest
}

func (c *CacheIngressManifest) SetSubnetManifest(cidr string, s *types.SubnetManifest) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for n := range c.manifests {

		if _, ok := c.manifests[n].Network[cidr]; !ok {
			c.manifests[n].Network = make(map[string]*types.SubnetManifest)
		}

		c.manifests[n].Network[cidr] = s
	}
}

func (c *CacheIngressManifest) SetRouteManifest(name string, s *types.RouteManifest) {
	c.lock.Lock()
	defer c.lock.Unlock()
	log.Debugf("set route manifest %s", name)

	if s.State == types.StateDestroyed {
		delete(c.routes, name)
	} else {
		c.routes[name] = s
	}

	for _, i := range c.ingress {
		if _, ok := c.manifests[i.SelfLink()]; ok {
			if _, ok := c.manifests[i.SelfLink()].Routes[name]; !ok {
				c.manifests[i.SelfLink()].Routes = make(map[string]*types.RouteManifest, 0)
			}
			c.manifests[i.SelfLink()].Routes[name] = s
		}
	}
}

func (c *CacheIngressManifest) SetEndpointManifest(addr string, s *types.EndpointManifest) {
	c.lock.Lock()
	defer c.lock.Unlock()

	log.Debugf("set endpoint manifest: %s > %#v", addr, s)

	for _, n := range c.manifests {
		if n.Endpoints == nil {
			n.Endpoints = make(map[string]*types.EndpointManifest, 0)
		}
		n.Endpoints[addr] = s
	}
}

func (c *CacheIngressManifest) SetIngress(ingress *types.Ingress) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.ingress[ingress.SelfLink()] = ingress
}

func (c *CacheIngressManifest) DelIngress(selflink string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.ingress, selflink)
	delete(c.manifests, selflink)
}

func (c *CacheIngressManifest) SetDiscovery(discovery *types.Discovery) {
	c.lock.Lock()
	defer c.lock.Unlock()

	dvc, ok := c.discovery[discovery.SelfLink()]

	if !ok {
		c.discovery[discovery.SelfLink()] = discovery
		c.SetResolvers()
		return
	}

	var update = false
	switch true {
	case dvc.Status.IP != discovery.Status.IP:
		update = true
		break
	case dvc.Status.Port != discovery.Status.Port:
		update = true
		break
	case dvc.Status.Ready != discovery.Status.Ready:
		update = true
		break
	}
	if update {
		c.discovery[discovery.SelfLink()] = discovery
		c.SetResolvers()
	}
	return
}

func (c *CacheIngressManifest) DelDiscovery(selflink string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.discovery, selflink)

	resolvers := make(map[string]*types.ResolverManifest, 0)

	for _, d := range c.discovery {
		if d.Status.Ready {
			resolvers[d.Status.IP] = &types.ResolverManifest{
				IP:   d.Status.IP,
				Port: d.Status.Port,
			}
		}
	}

	for _, n := range c.manifests {
		n.Meta.Discovery = resolvers
	}
}

func (c *CacheIngressManifest) SetResolvers() {
	resolvers := make(map[string]*types.ResolverManifest, 0)

	for _, d := range c.discovery {
		if d.Status.Ready {
			resolvers[d.Status.IP] = &types.ResolverManifest{
				IP:   d.Status.IP,
				Port: d.Status.Port,
			}
		}
	}

	for _, n := range c.manifests {
		n.Meta.Discovery = resolvers
	}
}

func (c *CacheIngressManifest) GetResolvers() map[string]*types.ResolverManifest {

	resolvers := make(map[string]*types.ResolverManifest, 0)

	for _, d := range c.discovery {
		if d.Status.Ready {
			resolvers[d.Status.IP] = &types.ResolverManifest{
				IP:   d.Status.IP,
				Port: d.Status.Port,
			}
		}
	}

	return resolvers
}

func (c *CacheIngressManifest) Get(ingress string) *types.IngressManifest {
	c.lock.Lock()
	defer c.lock.Unlock()
	if s, ok := c.manifests[ingress]; !ok {
		return nil
	} else {
		return s
	}
}

func (c *CacheIngressManifest) GetRoutes(ingress string) map[string]*types.RouteManifest {
	return c.routes
}

func (c *CacheIngressManifest) Flush(ingress string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.manifests[ingress] = new(types.IngressManifest)
}

func (c *CacheIngressManifest) Clear(ingress string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.manifests, ingress)
}

func NewCacheIngressManifest() *CacheIngressManifest {
	c := new(CacheIngressManifest)
	c.manifests = make(map[string]*types.IngressManifest, 0)
	c.discovery = make(map[string]*types.Discovery, 0)
	c.routes = make(map[string]*types.RouteManifest, 0)
	c.ingress = make(map[string]*types.Ingress)
	return c
}
