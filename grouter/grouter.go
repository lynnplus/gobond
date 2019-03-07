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
package grouter

import (
	"github.com/lynnsoft/gobond"
	"net/http"
	"strings"
)

type HttpRouter struct {

	initialized bool


	treeRoot map[string]*trieNode



	methodNotAllowedHandle gobond.HttpHandle
}

func New() gobond.GRouter{

	g:=&HttpRouter{}
	g.treeRoot= map[string]*trieNode{}
	g.initialized=false
	return g
}

func (r *HttpRouter) Initialized(){
	r.initialized=true
}


func (r *HttpRouter) SetRootHandler(handle gobond.HttpHandle) {
	panic("implement me")
}

func (r *HttpRouter) ServeFiles(path string, root http.FileSystem) {
	http.FileServer(root)
}


func (r *HttpRouter) GET(path string, handler gobond.HttpHandle) {
	r.addRoute(path,"GET",handler)
}

func (r *HttpRouter) POST(path string, handler gobond.HttpHandle) {
	r.addRoute(path,"POST",handler)
}

func (r *HttpRouter) HEAD(path string, handler gobond.HttpHandle) {
	r.addRoute(path,"HEAD",handler)
}

func (r *HttpRouter) PUT(path string, handler gobond.HttpHandle) {
	r.addRoute(path,"PUT",handler)
}

func (r *HttpRouter) CustomerHandle(method, path string, handler gobond.HttpHandle) {
	panic("implement me")
}

func (r *HttpRouter) HandlerFunc(method, path string, handler gobond.HttpHandle) {
	panic("implement me")
}


func (r *HttpRouter) ServeHTTP(ctx *gobond.WebContext){
	defer r.panicRecover(ctx)

	path:=ctx.Path()
	method:=ctx.Method()

	n:=r.treeRoot[method]
	if n==nil{
		if r.methodNotAllowedHandle!=nil{
			r.methodNotAllowedHandle(ctx)
		}
		return
	}

	paths:=strings.Split(path,"/")
	paths[0]="/"
	if path=="/"{
		paths=paths[:1]
	}
	found:=n.findNextNode(n,paths, func(h gobond.HttpHandle) bool {
		return h(ctx)
	})
	if !found{
		//404
	}
}


func (r *HttpRouter) panicRecover(ctx *gobond.WebContext) {
	if rcv := recover(); rcv != nil {

	}
}

func (r *HttpRouter) defaultRoot() *trieNode{
	t:=newTrieNode()
	t.path="/"
	return t
}


func (r *HttpRouter) addRoute(path string,method string,handle gobond.HttpHandle){
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if path=="/"{
		panic("use set '/' handler")
	}

	if r.treeRoot[method]==nil{
		r.treeRoot[method]=newTrieNode()
		r.treeRoot[method].childrenMap["/"]=r.defaultRoot()
	}

	rootNode:=r.treeRoot[method].childrenMap["/"]
	path=strings.TrimSuffix(path,"/")
	paths:=strings.Split(path,"/")
	r.treeRoot[method].insert(rootNode,paths,handle)
}



