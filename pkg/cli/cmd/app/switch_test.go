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

package app_test

//
//import (
//	"encoding/json"
//	n "github.com/lastbackend/lastbackend/pkg/api/app/views/v1"
//	"github.com/lastbackend/lastbackend/pkg/cli/cmd/app"
//	"github.com/lastbackend/lastbackend/pkg/cli/context"
//	storage "github.com/lastbackend/lastbackend/pkg/cli/storage/mock"
//	h "github.com/lastbackend/lastbackend/pkg/util/http"
//	"github.com/stretchr/testify/assert"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestSwitch(t *testing.T) {
//	const (
//		tName = "test name"
//		tDesc = "test description"
//	)
//
//	var (
//		err error
//		ctx = context.Mock()
//
//		data = n.App{
//			Meta: n.AppMeta{
//				Name:        tName,
//				Description: tDesc,
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
//
//		buff, err := json.Marshal(data)
//		assert.NoError(t, err)
//
//		w.WriteHeader(200)
//		_, err = w.Write(buff)
//		assert.NoError(t, err)
//	}))
//	defer server.Close()
//	//------------------------------------------------------------------------------------------
//
//	client, err := h.New(server.URL, &h.ReqOpts{})
//	assert.NoError(t, err)
//	ctx.SetHttpClient(client)
//
//	nspace, err := app.Switch(tName)
//	assert.NoError(t, err)
//	assert.Equal(t, tName, nspace.Meta.Name)
//	assert.Equal(t, tDesc, nspace.Meta.Description)
//}
