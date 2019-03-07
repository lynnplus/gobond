// Copyright 2019 gobond Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package gobond

import (
	"net/http"
	"sync"
)

type HttpHandle func(c *WebContext) (next bool)

type AppEngine struct {

	router GRouter
	contextPool sync.Pool
}


func NewApp() *AppEngine{
	app:=&AppEngine{}
	app.contextPool.New=app.newContext
	return app
}

func(a *AppEngine) Run(addr string,grouter GRouter) error{
	a.router=grouter
	a.router.Initialized()
	err := http.ListenAndServe(addr, a)
	return err
}



func (a *AppEngine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c:=a.contextPool.Get().(*WebContext)
	c.resetHttpRequest(w,req)
	a.handleHTTPRequest(c)
	a.contextPool.Put(c)
}


func (a *AppEngine) handleHTTPRequest(c *WebContext){
	a.router.ServeHTTP(c)
	return
}




func (a *AppEngine) newContext() interface{}{
	return &WebContext{app:a}
}