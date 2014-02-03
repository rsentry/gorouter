/*
Copyright (c) 2010/2011, Michael Beale
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

    * Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.

    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.

    * Neither the name of Michael Beale nor the names of contributors may be
      used to endorse or promote products derived from this software without
      specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL MICHAEL BEALE BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package router

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

//package vars
var router = new(Router).Init()

//package functions
//after
func After(urlquery string, handler ...func(w http.ResponseWriter, r *http.Request, v map[string]string)) {
	router.After(urlquery, handler)
}

//before
func Before(urlquery string, handler ...func(w http.ResponseWriter, r *http.Request, v map[string]string)) {
	router.Before(urlquery, handler)
}

//get
func Get(urlquery string, handler ...func(w http.ResponseWriter, r *http.Request, v map[string]string)) {
	router.Get(urlquery, handler)
}

//post
func Post(urlquery string, handler ...func(w http.ResponseWriter, r *http.Request, v map[string]string)) {
	router.Post(urlquery, handler)
}

//delete
func Delete(urlquery string, handler ...func(w http.ResponseWriter, r *http.Request, v map[string]string)) {
	router.Delete(urlquery, handler)
}

//put
func Put(urlquery string, handler ...func(w http.ResponseWriter, r *http.Request, v map[string]string)) {
	router.Put(urlquery, handler)
}

//run
func Run(address string) {
	http.ListenAndServe(address, router)
}

//overide errorhandler
func Handle404(new404handler func(w http.ResponseWriter, r *http.Request)) {
	router.error404Handler = new404handler
}

//set close flag to discontinue processing
func StopRequest(w http.ResponseWriter, r *http.Request) {
	r.Close = true
}

//handle 404 error
func Handle404Error(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 Not Found", http.StatusNotFound)
	StopRequest(w, r)
}

// Emits a 501 Not Implemented
func NotImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "501 Not Implemented", http.StatusNotImplemented)
	StopRequest(w, r)
}

// Emits a 201 Created with the URI for the new location
func Created(w http.ResponseWriter, r *http.Request, location string) {
	w.Header().Set("Location", location)
	http.Error(w, "201 Created", http.StatusCreated)
	StopRequest(w, r)
}

// Emits a 200 OK with a location. Used when after a PUT
func Updated(w http.ResponseWriter, r *http.Request, location string) {
	w.Header().Set("Location", location)
	http.Error(w, "200 OK", http.StatusOK)
	StopRequest(w, r)
}

// Emits a bad request with the specified instructions
func BadRequest(w http.ResponseWriter, r *http.Request, instructions string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(instructions))
	StopRequest(w, r)
}

// Emits a 204 No Content
func NoContent(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "204 No Content", http.StatusNoContent)
	StopRequest(w, r)
}

//structs
type Route struct {
	urlMap  string                                                              //a map of a url
	handler []func(w http.ResponseWriter, r *http.Request, v map[string]string) //handler function
}

type Router struct {
	post, get, put, delete, after, before []Route                                      //routes
	error404Handler                       func(w http.ResponseWriter, r *http.Request) //handler function
}

//end structs

//router functions
func (router *Router) Init() *Router {
	router = new(Router)
	router.error404Handler = Handle404Error
	return router
}

func (router *Router) AddRoute(method string, urlquery string, handler []func(w http.ResponseWriter, r *http.Request, v map[string]string)) {
	//add route to proper method
	var route Route
	route.urlMap = urlquery
	route.handler = handler
	switch method {
	case "get":
		router.get = append(router.get, route)
	case "post":
		router.post = append(router.post, route)
	case "delete":
		router.delete = append(router.delete, route)
	case "put":
		router.put = append(router.put, route)
	case "before":
		router.before = append(router.before, route)
	case "after":
		router.after = append(router.after, route)
	}
	return
}

func (router *Router) After(urlquery string, handler []func(w http.ResponseWriter, r *http.Request, v map[string]string)) {
	router.AddRoute("after", urlquery, handler)
}
func (router *Router) Before(urlquery string, handler []func(w http.ResponseWriter, r *http.Request, v map[string]string)) {
	router.AddRoute("before", urlquery, handler)
}

func (router *Router) Post(urlquery string, handler []func(w http.ResponseWriter, r *http.Request, v map[string]string)) {
	router.AddRoute("post", urlquery, handler)
}
func (router *Router) Get(urlquery string, handler []func(w http.ResponseWriter, r *http.Request, v map[string]string)) {
	router.AddRoute("get", urlquery, handler)
}
func (router *Router) Delete(urlquery string, handler []func(w http.ResponseWriter, r *http.Request, v map[string]string)) {
	router.AddRoute("delete", urlquery, handler)
}
func (router *Router) Put(urlquery string, handler []func(w http.ResponseWriter, r *http.Request, v map[string]string)) {
	router.AddRoute("Put", urlquery, handler)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//check for route in struct
	//maybe validate url?
	fmt.Printf("Route method: %s\n", r.URL.Path[1:])
	fmt.Printf("Route method: %s\n", r.Method)
	foundmatch := false
	switch r.Method {
	case "GET":
		foundmatch = findRouterMatch(router.get, r.URL.Path, w, r)
	case "POST":
		foundmatch = findRouterMatch(router.post, r.URL.Path, w, r)
	case "PUT":
		foundmatch = findRouterMatch(router.put, r.URL.Path, w, r)
	case "DELETE":
		foundmatch = findRouterMatch(router.delete, r.URL.Path, w, r)
	}
	if !foundmatch {
		//404
		router.error404Handler(w, r)
		fmt.Printf("No Route Found for url: %s \n", r.URL)
	}
}

//end router functions
func matchRoute(value Route, url string) (bool, map[string]string) {
	startlocation := strings.IndexAny(value.urlMap, "{*")
	valuemap := map[string]string{}
	//if url doesn't have any variables just do a straight match
	//else create regex
	switch {
	default:
		//create regex
		//init regex
		urlregex := "^"
		path := strings.Split(value.urlMap, "{")
		for _, fragment := range path {
			// varcheck := strings.Index(fragment, ":")
			varcheck := strings.IndexAny(fragment, ":*")
			switch {
			default:
				//no special characters found
				urlregex += fragment
			case varcheck >= 0:
				if varcheck := strings.Index(fragment, ":"); varcheck >= 0 {
					//it is a variable
					//find text after closing bracket
					bracketfragment := strings.Split(fragment, "}")
					//remove : to get variable name
					variablename := strings.Replace(bracketfragment[0], ":", "", -1)
					//remove extra from string
					rx := regexp.MustCompile(urlregex)
					varstring := rx.ReplaceAllString(url, "")
					rx = regexp.MustCompile("[a-zA-Z0-9_-]*")
					//get filtered string
					matched := rx.FindStringSubmatch(varstring)
					switch {
					default:
						valuemap[variablename] = "unknown"
					case matched != nil:
						valuemap[variablename] = matched[0]
					}
					urlregex += "[^/]"
					//get leftover text
					if len(bracketfragment) > 1 && bracketfragment[1] != "" {
						//replace . with \. 
						//todo:  look into other characters that need to be escaped
						urlregex += strings.Replace(bracketfragment[1], ".", "\\.", -1)
					}
				} else {
					//need to handle * and  *.*
					urlregex += "."
				}
			}
		}
		urlregex += "*"
		rx := regexp.MustCompile(urlregex)
		matched := rx.ReplaceAllString(url, "")
		//if string == "" that means there is a match
		if matched == "" {
			return true, valuemap
		}
	case startlocation < 0:
		//this is when no "{" is found so just does a regular string compare
		if url == value.urlMap {
			return true, nil
		}
	}
	return false, nil
}
func findRouterMatch(routes []Route, url string, w http.ResponseWriter, r *http.Request) bool {
	router.findBefore(url, w, r)
	if !r.Close {
		for _, value := range routes {
			//todo: add pass functionality
			if matched, valuemap := matchRoute(value, url); matched {
				fmt.Printf("Route executing %s\n", value.urlMap)
				for _, f := range value.handler {
					if !r.Close {
						f(w, r, valuemap)
					}
				}
				return true
			}
		}
	}
	router.findAfter(url, w, r)
	return false
}

func (router *Router) findAfter(url string, w http.ResponseWriter, r *http.Request) {
	for _, value := range router.after {
		fmt.Printf("Route executing %s against url %s\n", value.urlMap, url)
		if value.urlMap == "" {
			//execute for all requests
			valuemap := map[string]string{}
			for _, f := range value.handler {
				if !r.Close {
					f(w, r, valuemap)
				}
			}
		} else if matched, valuemap := matchRoute(value, url); matched {
			fmt.Printf("Route executing %s\n", value.urlMap)
			for _, f := range value.handler {
				if !r.Close {
					f(w, r, valuemap)
				}
			}
		}
	}
}
func (router *Router) findBefore(url string, w http.ResponseWriter, r *http.Request) {
	// for _, v := range router.before {
	// 	v.handler(w, r)
	// }
}
