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
)

type trieNode struct {
	path      string
	wildChild bool
	maxParams uint8
	childrenMap  map[string]*trieNode
	handle    gobond.HttpHandle
	priority  uint32
}

type callback func(h gobond.HttpHandle) bool


func newTrieNode() *trieNode{
	return &trieNode{
		childrenMap: make(map[string]*trieNode),
	}
}


func (t *trieNode) insert(parent *trieNode,paths []string,handle gobond.HttpHandle){
	if len(paths)==0{
		panic("err")
		return
	}

	first:=paths[0]
	if first==""{
		t.insert(parent,paths[1:],handle)
		return
	}

	firstNode:=parent.find(first)
	if firstNode==nil{
		t.insertBranch(parent,paths,handle)
		return
	}

	if len(paths)==1{
		firstNode.path=first
		firstNode.handle=handle
		return
	}
	rest:=paths[1:]
	t.insert(firstNode,rest,handle)
}

func (t *trieNode) insertBranch(parent *trieNode,paths []string,handle gobond.HttpHandle){
	if len(paths)==0{
		panic("error")
	}

	first:=paths[0]

	n:=newTrieNode()
	parent.childrenMap[first]=n

	if len(paths)==1{
		n.handle=handle
		return
	}
	rest:=paths[1:]
	t.insertBranch(n,rest,handle)
}

func (t *trieNode) emptyChilds() bool {
	return len(t.childrenMap) == 0
}


func (t *trieNode) find(path string) *trieNode{
	return t.childrenMap[path]
}

func (t *trieNode) findNextNode(parent *trieNode,paths []string,fn callback) bool{
	if len(paths)==0{
		return false
	}

	first:=paths[0]
	firstNode:=parent.find(first)
	if firstNode==nil{
		return false
	}
	if firstNode.handle==nil{
		rest:=paths[1:]
		return t.findNextNode(firstNode,rest,fn)
	}

	next:=fn(firstNode.handle)
	if !next{
		return false
	}

	if len(paths)==1{

		return true
	}
	
	if firstNode.emptyChilds(){
		return false
	}
	
	rest:=paths[1:]
	return t.findNextNode(firstNode,rest,fn)
}
