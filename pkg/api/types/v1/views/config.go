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

package views

import (
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"time"
)

// swagger:model views_secret
type Config struct {
	Meta ConfigMeta `json:"meta"`
	Spec ConfigSpec `json:"spec"`
}

type ConfigSpec struct {
	Type string           `json:"type"`
	Data []ConfigSpecData `json:"data"`
}

type ConfigSpecData struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
	File  string `json:"file,omitempty"`
	Data  []byte `json:"data,omitempty"`
}

// swagger:model views_secret_meta
type ConfigMeta struct {
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Kind      string    `json:"kind"`
	SelfLink  string    `json:"self_link"`
	Updated   time.Time `json:"updated"`
	Created   time.Time `json:"created"`
}

// swagger:ignore
type ConfigMap map[string]*Config

// swagger:model views_secret_list
type ConfigList []*Config

func (s *Config) Decode() *types.Config {

	o := new(types.Config)
	o.Meta.Name = s.Meta.Name
	o.Meta.Namespace = s.Meta.Namespace
	o.Meta.Kind = s.Meta.Kind
	o.Meta.SelfLink = s.Meta.SelfLink
	o.Meta.Updated = s.Meta.Updated
	o.Meta.Created = s.Meta.Created

	o.Spec.Type = s.Spec.Type
	o.Spec.Data = make([]*types.ConfigSpecData, 0)

	for _, d := range s.Spec.Data {
		o.Spec.Data = append(o.Spec.Data, &types.ConfigSpecData{
			Key: d.Key,
			File: d.File,
			Value: d.Value,
			Data: d.Data,
		})
	}

	return o
}
