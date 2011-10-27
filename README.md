##gorouter

A simple router for go

##Overview
This project is a simple lightweight router for go. Some of the features are:

* Regex routing
* Variable Mapping

##Installation
You need a working installation of go obviously.

The easiest way:
goinstall github.com/rsentry/gorouter

Or you can download and do make && make install

##Example

	package main
	
	import (
			"http"
			"fmt"
			"github.com/rsentry/gorouter"
	)
       
	func itemshandler(w http.ResponseWriter, r *http.Request, v map[string] string){
		fmt.Fprintf(w, "Hi there this is items<br/>")
		for vkey,vvalue := range v {
			fmt.Fprintf(w,"varname:%s , value:%s<br/>", vkey,vvalue)
		}
	}

	
	func saleshandler(w http.ResponseWriter, r *http.Request , v map[string] string){
		fmt.Fprintf(w, "Hi there this is sales")
	}
	
	func purchaseinghandler(w http.ResponseWriter, r *http.Request , v map[string] string){
		fmt.Fprintf(w, "Hi there this is purchase orders")
	}
	
	func errorHandler(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "This is an overide handler")
	}
	func main(){
		router.Get("/v1/items/{:id}.{:type}", itemshandler)
		router.Get("/v1/items.{:type}", itemshandler)
		router.Get("/v1/purchase_orders", purchaseinghandler)
		router.Get("sales{*}", purchaseinghandler)
		router.Handle404(errorHandler)
		router.Run(":8080")
	}

To run the app do:
	8g router.go && 8g hello.go && 8l -o hello hello.8 && ./hello

then point your browser to http://localhost:8080/items/1.json

##Creating routes
anytime you want a variable in the url, surround it with {:}.  This will extract the variable name and the url data and place the information in the variable map.
If there are not any variables, then the url must match your route exactly. You can also override the 404 handler to catch 404 errors. See the example section for more detail.

##Multiple Handlers
You can add mutliple handlers to a specific route.  For instance:
	router.Get("/v1/items", checkifauthhandler,itemsgethandler)
The handlers will be run in the sequence specified.  If one of your handlers wants to terminate processing the request, call the function StopRequest(w , r) and the remaining handlers will not be called.  If you used variables in the url, the variables will be in the variable map. Also, you can call the following functions to stop the request:
	
	func Handle404Error(w http.ResponseWriter, r *http.Request)
	func NotImplemented(w http.ResponseWriter, r *http.Request) 
	func Created(w http.ResponseWriter, r *http.Request, location string) 
	func Updated(w http.ResponseWriter, r *http.Request, location string) 
	func BadRequest(w http.ResponseWriter, r *http.Request, instructions string) 
	func NoContent(w http.ResponseWriter, r *http.Request) 

##Future
This is a lightweight router so there are not many features planned but here is a list of a few that I would like to implement:

* Custom Logging
* Hooks called before and after handler
* Write Tests
* Write godocs
* Static File Logic
* Route Passing and Calling
* File handling
* Specifying multiple methods for each route
* Open to suggestions

##About
This should be considered very alpha and not recommended for production use.
gorouter was written by Michael Beale
