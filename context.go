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
	"encoding/json"
	"net/http"
)

type WebContext struct {

	app *AppEngine


	req *http.Request
	resp http.ResponseWriter
}


func (c *WebContext) Path() string{
	return c.req.URL.Path
}

func (c *WebContext) Method() string{
	return c.req.Method
}

func (c *WebContext) resetHttpRequest(w http.ResponseWriter, req *http.Request){
	c.req=req
	c.resp=w
}

func (c *WebContext) resetRpcRequest(){



}

func (c *WebContext) ResponseJson(data interface{}) error{
	return c.handleResponse(http.StatusOK,"JSON",data)
}

func (c *WebContext) handleResponse(code int,dataType string,data interface{}) error{
	s,e:=json.Marshal(data)
	if e!=nil{
		return e
	}
	c.resp.WriteHeader(code)
	header:=c.resp.Header()
	if val:=header["Content-Type"];len(val)==0{
		header["Content-Type"]=[]string{"application/json; charset=utf-8"}
	}
	_,e=c.resp.Write(s)
	if e!=nil{
		return e
	}
	return nil
}


