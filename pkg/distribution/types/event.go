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

package types

const (
	EventActionCreate = "create"
	EventActionUpdate = "update"
	EventActionDelete = "delete"
)

type event struct {
	Action string
	Name   string
}

type Event struct {
	event
	Data interface{}
}

type NamespaceEvent struct {
	event
	Data *Namespace
}

type ClusterEvent struct {
	event
	Data *Cluster
}

type ServiceEvent struct {
	event
	Data *Service
}

type IngresEvent struct {
	event
	Data *Ingress
}

type EndpointEvent struct {
	event
	Data *Endpoint
}
