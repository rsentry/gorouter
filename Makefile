include $(GOROOT)/src/Make.inc

TARG=gorouter
GOFMT=gofmt -s -spaces=true -tabindent=false -tabwidth=4

GOFILES=\
    router.go\

include $(GOROOT)/src/Make.pkg
