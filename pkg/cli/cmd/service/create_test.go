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

package service_test

//
//import (
//	"encoding/json"
//	n "github.com/lastbackend/lastbackend/pkg/api/app/views/v1"
//	sm "github.com/lastbackend/lastbackend/pkg/api/service/views/v1"
//	"github.com/lastbackend/lastbackend/pkg/cli/cmd/service"
//	"github.com/lastbackend/lastbackend/pkg/cli/context"
//	storage "github.com/lastbackend/lastbackend/pkg/cli/storage/mock"
//	h "github.com/lastbackend/lastbackend/pkg/util/http"
//	"github.com/stretchr/testify/assert"
//	"io/ioutil"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestCreate(t *testing.T) {
//
//	const (
//		nName = "nspace name"
//
//		sName = "service name"
//		sDesc = "service desc"
//	)
//
//	var (
//		err error
//		ctx = context.Mock()
//
//		ns = &n.App{
//			Meta: n.AppMeta{
//				Name: nName,
//			},
//		}
//	)
//
//	strg, err := storage.Get()
//	assert.NoError(t, err)
//	ctx.SetStorage(strg)
//	defer strg.App().Remove()
//
//	//------------------------------------------------------------------------------------------
//	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		body, err := ioutil.ReadAll(r.Body)
//		assert.NoError(t, err)
//
//		var d = struct {
//			Name string `json:"name,omitempty"`
//		}{}
//
//		err = json.Unmarshal(body, &d)
//		assert.NoError(t, err)
//
//		assert.Equal(t, sName, d.Name)
//
//		nspaceJSON, err := json.Marshal(sm.Service{
//			Meta: sm.ServiceMeta{
//				Name:        sName,
//				Description: sDesc,
//			},
//		})
//		assert.NoError(t, err)
//
//		w.WriteHeader(200)
//		_, err = w.Write(nspaceJSON)
//		assert.NoError(t, err)
//	}))
//	defer server.Close()
//	//------------------------------------------------------------------------------------------
//
//	strg.App().Save(ns)
//	assert.NoError(t, err)
//
//	client, err := h.New(server.URL, &h.ReqOpts{})
//	assert.NoError(t, err)
//	ctx.SetHttpClient(client)
//
//	err = service.Create(sName, "redis", "", "", nil)
//	assert.NoError(t, err)
//}
