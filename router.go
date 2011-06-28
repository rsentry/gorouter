package router

import (
		"fmt"
		"http"
		"strings"
		"regexp"
       )
//package vars
var	router = new(Router).Init()


//package functions
//get
func Get(urlquery string, handler func(w http.ResponseWriter, r *http.Request, v map[string] string)){
	router.Get(urlquery , handler)
}
//post
func Post(urlquery string, handler func(w http.ResponseWriter, r *http.Request, v map[string] string)){
	router.Post(urlquery , handler)
}
//delete
func Delete(urlquery string, handler func(w http.ResponseWriter, r *http.Request, v map[string] string)){
	router.Delete(urlquery , handler)
}
//put
func Put(urlquery string, handler func(w http.ResponseWriter, r *http.Request, v map[string] string)){
	router.Put(urlquery , handler)
}
//run
func Run(address string){
	http.ListenAndServe(address, router)
}
//overide errorhandler
func Handle404(new404handler func (w http.ResponseWriter, r *http.Request)){
	router.error404Handler = new404handler
}

func handle404Error(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Unable to locate resource")
}
//structs
type Route struct {
	urlMap string //a map of a url
	handler func(w http.ResponseWriter, r *http.Request, v map[string] string) //handler function
}

type Router struct {
	post,get,put,delete []Route //routes
	error404Handler func(w http.ResponseWriter, r *http.Request) //handler function
}


//end structs

//router functions
func (router *Router) Init() *Router{
	router = new(Router)
	router.error404Handler = handle404Error
	return router
}

func (router *Router) AddRoute(method string, urlquery string, handler func(w http.ResponseWriter, r *http.Request, v map[string] string)) {
	//add route to proper method
	var route Route
	route.urlMap = urlquery
	route.handler = handler
	switch method {
	case "get":
		router.get = append(router.get,route)
	case "post":
		router.post = append(router.post,route)
	case "delete":
		router.delete = append(router.delete,route)
	case "put":
		router.put = append(router.put,route)
	}
	return
}

func (router *Router) Post(urlquery string,handler func(w http.ResponseWriter, r *http.Request, v map[string] string)){
	router.AddRoute("post", urlquery, handler)
}
func (router *Router) Get(urlquery string, handler func(w http.ResponseWriter, r *http.Request, v map[string] string)){
	router.AddRoute("get", urlquery, handler)
}
func (router *Router) Delete(urlquery string, handler func(w http.ResponseWriter, r *http.Request, v map[string] string)){
	router.AddRoute("delete", urlquery, handler)
}
func (router *Router) Put(urlquery string, handler func(w http.ResponseWriter, r *http.Request, v map[string] string)){
	router.AddRoute("Put", urlquery, handler)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request){ 
	//check for route in struct
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
		router.error404Handler(w,r)
		fmt.Printf("No Route Found for url: %s \n", r.RawURL)
	}
}
//end router functions

func findRouterMatch(routes []Route, url string,w http.ResponseWriter, r *http.Request ) bool{
	//todo:  add hooks for before and after call
	for _ ,value := range routes {
		startlocation := strings.Index(value.urlMap, "{")
		valuemap := map[string] string{} 
		//if url doesn't have any variables just do a straight match
		//else create regex
		switch{
		default:
			//create regex
			//init regex
			urlregex := "^"
			path := strings.Split(value.urlMap,"{",-1)
			for _ , fragment := range path{
				varcheck := strings.Index(fragment,":")
				switch{
				default:
					//no special characters found
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
				value.handler(w , r, valuemap)
				return true
			}
		case startlocation < 0:
			//this is when no "{" is found so just does a regular string compare
			if url == value.urlMap {
				//set defaults(if not set) for offset, limit, select
				fmt.Printf("Route executing %s\n", value.urlMap)
				value.handler(w , r, valuemap)
				return true
			}
		}
	}
	return false	
}
