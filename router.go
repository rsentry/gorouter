package router

import (
		"fmt"
		"http"
		"strings"
		"regexp"
       )
//package vars
var	router = new(Router).Init(100)


//package functions
//get
func Get(urlquery string, handler func(w http.ResponseWriter, r *http.Request)){
	router.Get(urlquery , handler)
}
//post
func Post(urlquery string, handler func(w http.ResponseWriter, r *http.Request)){
	router.Post(urlquery , handler)
}
//delete
func Delete(urlquery string, handler func(w http.ResponseWriter, r *http.Request)){
	router.Delete(urlquery , handler)
}
//put
func Put(urlquery string, handler func(w http.ResponseWriter, r *http.Request)){
	router.Put(urlquery , handler)
}
//run
func Run(address string){
	http.ListenAndServe(address, router)
}
//structs
type Route struct {
	urlMap string //a map of a url
	handler func(w http.ResponseWriter, r *http.Request) //handler function
}

type Router struct {
	post,get,put,delete []Route
}


//end structs

//router functions
func (router *Router) Init(routecap int) *Router{
	router = new(Router)
	router.post = make([]Route,0,routecap)
	router.put = make([]Route,0,routecap)
	router.delete = make([]Route,0,routecap)
	router.get = make([]Route,0,routecap)
	return router
}

func (router *Router) AddRoute(method string, urlquery string, handler func(w http.ResponseWriter, r *http.Request)) {
	//add route to proper method
	var route Route
	route.urlMap = urlquery
	route.handler = handler
	switch method {
	case "get":
		n := len(router.get)
		router.get = router.get[0:n+1]
		router.get[n] = route
	case "post":
		n := len(router.post)
		router.post = router.post[0:n+1]
		router.post[n] = route
	case "delete":
		n := len(router.delete)
		router.delete = router.delete[0:n+1]
		router.delete[n] = route
	case "put":
		n := len(router.put)
		router.put = router.put[0:n+1]
		router.put[n] = route
	}
	return
}

func (router *Router) Post(urlquery string,handler func(w http.ResponseWriter, r *http.Request)){
	router.AddRoute("post", urlquery, handler)
}
func (router *Router) Get(urlquery string, handler func(w http.ResponseWriter, r *http.Request)){
	router.AddRoute("get", urlquery, handler)
}
func (router *Router) Delete(urlquery string, handler func(w http.ResponseWriter, r *http.Request)){
	router.AddRoute("delete", urlquery, handler)
}
func (router *Router) Put(urlquery string, handler func(w http.ResponseWriter, r *http.Request)){
	router.AddRoute("Put", urlquery, handler)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request){ 
	//check for route in struct
	path := strings.Split(r.URL.Path,"/",-1)
	for key , value := range path{
		fmt.Printf("Route url key %s: %s\n", key, value)
	}
	fmt.Printf("Route method: %s\n", r.URL.Path[1:])
	fmt.Printf("Route method: %s\n", r.Method)
	foundmatch := false
	switch r.Method{
	case "GET":
		foundmatch = findRouterMatch(router.get, r.URL.Path, w, r)
	case "POST":
		foundmatch = findRouterMatch(router.post, r.URL.Path, w , r)
	case "PUT":
		foundmatch = findRouterMatch(router.put, r.URL.Path, w , r)
	case "DELETE":
		foundmatch = findRouterMatch(router.delete, r.URL.Path, w , r)
	}
	if(!foundmatch){
		//404
		fmt.Fprintf(w, "No Route Found for url: %s \n", r.RawURL)
	}
}
//end router functions

func findRouterMatch(routes []Route, url string,w http.ResponseWriter, r *http.Request ) bool{
	for _ ,value := range routes {
		startlocation := strings.Index(value.urlMap, "{")
		valuemap := map[string] string{} 
		//if url doesn't have any variables just do a straight match
		//else create regex
		switch{
		default:
			//create regex
			urlregex := "^"
			path := strings.Split(value.urlMap,"{",-1)
			for _ , fragment := range path{
				varcheck := strings.Index(fragment,":")
				switch{
				default:
					urlregex += fragment
				case varcheck >= 0:
					//it is a variable
					//find text after closing bracket
					bracketfragment := strings.Split(fragment,"}",-1)
					//remove : to get variable name
					variablename := strings.Replace(bracketfragment[0],":","",-1)
					//remove extra from string
					rx := regexp.MustCompile(urlregex)
					varstring := rx.ReplaceAllString(url,"")
					rx = regexp.MustCompile("[a-zA-Z0-9_]*")
					matched := rx.FindStringSubmatch(varstring)
					switch {
					default:
						valuemap[variablename] = "unknown"
					case matched != nil:
						valuemap[variablename] = matched[0]
					}
						
					// urlregex += "[a-zA-Z0-9_]"
					urlregex += "[^/]"
					//get leftover text
					if len(bracketfragment) > 1 && bracketfragment[1] != "" {
						//replace . with \. 
						urlregex += strings.Replace(bracketfragment[1], ".", "\\.",-1)
					}
				//need to handle * and  *.*
				}
			}
			urlregex += "*"
			rx := regexp.MustCompile(urlregex)
			matched := rx.ReplaceAllString(url, "")
			//if string == "" that means there is a match
			if matched == ""{
				fmt.Printf("Route executing %s\n", value.urlMap)
				fmt.Printf("Regex match regex:%s , matched string:%s\n", rx.String(), url)
				for vkey,vvalue := range valuemap {
					fmt.Printf("varname:%s , value:%s\n", vkey,vvalue)
				}
				value.handler(w , r)
				return true
			}
		case startlocation < 0:
			if url == value.urlMap {
				//create variable map
				//set defaults(if not set) for offset, limit, select
				fmt.Printf("Route executing %s\n", value.urlMap)
				value.handler(w , r)
				return true
			}
		}
	}
	return false	
}
