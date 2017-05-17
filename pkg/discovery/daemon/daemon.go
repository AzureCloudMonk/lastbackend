//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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

package daemon

import (
	_cfg "github.com/lastbackend/lastbackend/pkg/common/config"
	"github.com/lastbackend/lastbackend/pkg/discovery/cache"
	"github.com/lastbackend/lastbackend/pkg/discovery/config"
	"github.com/lastbackend/lastbackend/pkg/discovery/context"
	"github.com/lastbackend/lastbackend/pkg/discovery/domain"
	"github.com/lastbackend/lastbackend/pkg/logger"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"os"
	"os/signal"
	"syscall"
)

func Daemon(_cfg *_cfg.Config) {

	var (
		ctx  = context.Get()
		cfg  = config.Set(_cfg)
		log  = logger.New("Discovery", *cfg.LogLevel)
		sigs = make(chan os.Signal)
		done = make(chan bool, 1)
	)

	log.Info("Start service discovery")

	ctx.SetConfig(cfg)
	ctx.SetLogger(log)
	ctx.SetCache(cache.New())

	stg, err := storage.Get(cfg.GetEtcdDB())
	if err != nil {
		panic(err)
	}
	ctx.SetStorage(stg)

	sd, err := Listen(*cfg.DiscoveryServer.Port)
	if err != nil {
		log.Warnf("Start discovery server error: %s", err.Error())
	}

	go func() {
		if err := domain.Watch(); err != nil {
			log.Warnf("Watch domain error: %s", err.Error())
		}
	}()

	// Handle SIGINT and SIGTERM.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-sigs:
				done <- true
				return
			}
		}
	}()

	<-done

	sd.Shutdown()

	log.Info("Handle SIGINT and SIGTERM.")
}
