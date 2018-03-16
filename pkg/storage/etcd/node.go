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

package etcd

import (
	"context"
	"errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"regexp"
	"time"
)

const (
	nodeStorage = "node"
	timeout     = 15
)

// Node Service type for interface in interfaces folder
type NodeStorage struct {
	storage.Node
}

func (s *NodeStorage) List(ctx context.Context) ([]*types.Node, error) {

	log.V(logLevel).Debugf("Storage: Node: get list nodes")

	const filter = `\b.+` + nodeStorage + `\/(.+)\/(?:meta|state|alive)\b`

	nodes := make([]*types.Node, 0)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: Node: create client err: %s", err.Error())
		return nil, err
	}
	defer destroy()

	key := keyCreate(nodeStorage)
	if err := client.List(ctx, key, filter, &nodes); err != nil {
		log.V(logLevel).Errorf("Storage: Node: get nodes list err: %s", err.Error())
		return nil, err
	}

	log.V(logLevel).Debugf("Storage: Node: get nodes list result: %d", len(nodes))

	return nodes, nil
}

func (s *NodeStorage) Get(ctx context.Context, id string) (*types.Node, error) {

	log.V(logLevel).Debugf("Storage: Node: get by id: %s", id)

	if len(id) == 0 {
		err := errors.New("id can not be empty")
		log.V(logLevel).Errorf("Storage: Node: get node err: %s", err.Error())
		return nil, err
	}

	const filter = `\b.+` + nodeStorage + `\/.+\/(?:meta|state|alive)\b`

	client, destroy, err := getClient(ctx)
	if err != nil {
		return nil, err
	}
	defer destroy()

	node := new(types.Node)
	key := keyCreate(nodeStorage, id)
	if err := client.Map(ctx, key, filter, node); err != nil {
		if err.Error() == store.ErrKeyNotFound {
			return nil, nil
		}
		log.V(logLevel).Errorf("Storage: Node: create client err: %s", err.Error())
		return nil, err
	}

	if node.Meta.Name == "" {
		return nil, errors.New(store.ErrKeyNotFound)
	}

	node.Spec.Pods = make(map[string]*types.Pod)
	keySpec := keyCreate(nodeStorage, id, "spec", "pods")
	if err := client.Map(ctx, keySpec, "", node.Spec.Pods); err != nil {
		// Return node if pods does not exists
		if err.Error() == store.ErrKeyNotFound {
			return node, nil
		}
		log.V(logLevel).Errorf("Storage: Node: get node err: %s", err.Error())
		return nil, err
	}

	return node, nil
}

func (s *NodeStorage) Insert(ctx context.Context, node *types.Node) error {

	log.V(logLevel).Debugf("Storage: Node: insert node: %#v", node)

	if node == nil {
		err := errors.New("node can not be nil")
		log.V(logLevel).Errorf("Storage: Node: insert node err: %s", err.Error())
		return err
	}

	node.Meta.Created = time.Now()
	node.Meta.Updated = time.Now()

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: Node: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyMeta := keyCreate(nodeStorage, node.Meta.Name, "meta")
	if err := tx.Create(keyMeta, &node.Meta, 0); err != nil {
		log.V(logLevel).Errorf("Storage: Node: create meta err: %s", err.Error())
		return err
	}

	keyState := keyCreate(nodeStorage, node.Meta.Name, "state")
	if err := tx.Create(keyState, &node.State, 0); err != nil {
		log.V(logLevel).Errorf("Storage: Node: create state err: %s", err.Error())
		return err
	}

	keyAvailable := keyCreate(nodeStorage, node.Meta.Name, "alive")
	if err := tx.Create(keyAvailable, true, timeout); err != nil {
		log.V(logLevel).Errorf("Storage: Node: create alive err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("Storage: Node: commit transaction err: %s", err.Error())
		return err
	}

	return nil
}

func (s *NodeStorage) Update(ctx context.Context, node *types.Node) error {

	log.V(logLevel).Debugf("Storage: Node: update node: %#v", node)

	if node == nil {
		err := errors.New("node can not be nil")
		log.V(logLevel).Errorf("Storage: Node: update node err: %s", err.Error())
		return err
	}

	node.Meta.Updated = time.Now()

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: Node: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyMeta := keyCreate(nodeStorage, node.Meta.Name, "meta")
	if err := tx.Update(keyMeta, &node.Meta, 0); err != nil {
		log.V(logLevel).Errorf("Storage: Node: update node meta err: %s", err.Error())
		return err
	}

	keyState := keyCreate(nodeStorage, node.Meta.Name, "state")
	if err := tx.Update(keyState, &node.State, 0); err != nil {
		log.V(logLevel).Errorf("Storage: Node: update node state err: %s", err.Error())
		return err
	}

	keyAvailable := keyCreate(nodeStorage, node.Meta.Name, "alive")
	if err := tx.Upsert(keyAvailable, true, timeout); err != nil {
		log.V(logLevel).Errorf("Storage: Node: upsert node alive err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("Storage: Node: commit transaction err: %s", err.Error())
		return err
	}

	return nil
}

func (s *NodeStorage) InsertPod(ctx context.Context, meta *types.NodeMeta, pod *types.Pod) error {

	log.V(logLevel).Debugf("Storage: Node: insert pod in node: %#v", pod)

	if meta == nil {
		err := errors.New("meta can not be nil")
		log.V(logLevel).Errorf("Storage: Node: insert pod in node err: %s", err.Error())
		return err
	}

	if pod == nil {
		err := errors.New("pod can not be nil")
		log.V(logLevel).Errorf("Storage: Node: insert node in pod err: %s", err.Error())
		return err
	}

	meta.Updated = time.Now()

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: Node: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)
	keyMeta := keyCreate(nodeStorage, meta.Name, "meta")
	if err := tx.Update(keyMeta, meta, 0); err != nil {
		log.V(logLevel).Errorf("Storage: Node: update node meta err: %s", err.Error())
		return err
	}

	keyPod := keyCreate(nodeStorage, meta.Name, "spec", "pods", pod.Meta.Name)
	if err := tx.Create(keyPod, pod, 0); err != nil {
		log.V(logLevel).Errorf("Storage: Node: create pod for node err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("Storage: Node: commit transaction err: %s", err.Error())
		return err
	}

	return nil
}

func (s *NodeStorage) UpdatePod(ctx context.Context, meta *types.NodeMeta, pod *types.Pod) error {

	log.V(logLevel).Debugf("Storage: Node: update pod in node: %#v", pod)

	if meta == nil {
		err := errors.New("meta can not be nil")
		log.V(logLevel).Errorf("Storage: Node: update pod in node err: %s", err.Error())
		return err
	}

	if pod == nil {
		err := errors.New("pod can not be nil")
		log.V(logLevel).Errorf("Storage: Node: update pod in node err: %s", err.Error())
		return err
	}

	meta.Updated = time.Now()

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: Node: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)
	keyMeta := keyCreate(nodeStorage, meta.Name, "meta")
	if err := tx.Update(keyMeta, meta, 0); err != nil {
		log.V(logLevel).Errorf("Storage: Node: update node meta err: %s", err.Error())
		return err
	}

	keyPod := keyCreate(nodeStorage, meta.Name, "spec", "pods", pod.Meta.Name)
	if err := tx.Update(keyPod, pod, 0); err != nil {
		log.V(logLevel).Errorf("Storage: Node: update node spec pods err: %s", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("Storage: Node: commit transaction err: %s", err.Error())
		return err
	}

	return nil
}

func (s *NodeStorage) RemovePod(ctx context.Context, meta *types.NodeMeta, pod *types.Pod) error {

	log.V(logLevel).Debugf("Storage: Node: remove pod from node: %#v", pod)

	if meta == nil {
		err := errors.New("meta can not be nil")
		log.V(logLevel).Errorf("Storage: Node: remove pod from node err: %s", err.Error())
		return err
	}

	if pod == nil {
		err := errors.New("pod can not be nil")
		log.V(logLevel).Errorf("Storage: Node: remove pod from node err: %s", err.Error())
		return err
	}

	meta.Updated = time.Now()

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: Node: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	tx := client.Begin(ctx)

	keyMeta := keyCreate(nodeStorage, meta.Name, "meta")
	if err := tx.Update(keyMeta, meta, 0); err != nil {
		log.V(logLevel).Errorf("Storage: Node: update node meta err: %s", err.Error())
		return err
	}

	keyPod := keyCreate(nodeStorage, meta.Name, "spec", "pods", pod.Meta.Name)
	tx.Delete(keyPod)

	if err := tx.Commit(); err != nil {
		log.V(logLevel).Errorf("Storage: Node: commit transaction err: %s", err.Error())
		return err
	}

	return nil
}

func (s *NodeStorage) Remove(ctx context.Context, name string) error {

	log.V(logLevel).Debugf("Storage: Node: remove node: %#v", name)

	if name == "" {
		err := errors.New("node can not be nil")
		log.V(logLevel).Errorf("Storage: Node: remove node err: %s", err.Error())
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: Node: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(nodeStorage, name)
	if err := client.DeleteDir(ctx, key); err != nil {
		log.V(logLevel).Errorf("Storage: Node: remove node err: %s", err.Error())
		return err
	}

	return nil
}

func (s *NodeStorage) Watch(ctx context.Context, node chan *types.Node) error {

	log.V(logLevel).Debug("Storage: Node: watch node")

	const filter = `\b.+` + nodeStorage + `\/(.+)\/alive\b`

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: Node: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	r, _ := regexp.Compile(filter)
	key := keyCreate(nodeStorage)
	cb := func(action, key string, _ []byte) {
		keys := r.FindStringSubmatch(key)
		if len(keys) < 2 {
			return
		}

		n, _ := s.Get(ctx, keys[1])
		if n == nil {
			return
		}

		// TODO: check previous node alive state to prevent multi calls
		if action == "PUT" {
			n.Alive = true
			node <- n
			return
		}

		if action == "DELETE" {
			n.Alive = false
			node <- n
			return
		}

		return
	}

	if err := client.Watch(ctx, key, filter, cb); err != nil {
		log.V(logLevel).Errorf("Storage: Node: watch node err: %s", err.Error())
		return err
	}

	return nil
}

func newNodeStorage() *NodeStorage {
	s := new(NodeStorage)
	return s
}
