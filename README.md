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
		router.Handle404(errorHandler)
		router.Run(":8080")
	}

To run the app do:
	8g router.go && 8g hello.go && 8l -o hello hello.8 && ./hello

then point your browser to http://localhost:8080/items/1.json

##Creating routes
anytime you want a variable in the url, surround it with {:}.  This will extract the variable name and the url data and place the information in the variable map.
If there are not any variables, then the url must match your route exactly. You can also override the 404 handler to catch 404 errors. See the example section for more detail.

##Future
This is a lightweight router so there are not many features planned but here is a list of a few that I would like to implement:

* Custom Logging
* Hooks called before and after handler
* More complex url handling
* Open to suggestions

##About

gorouter was written by Michael Beale
