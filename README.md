gae-go-starter
==============

web application starter template for GAE/go

#set up
##1. install go
if you use Mac OS, "brew install go" is the easiest way.

set PATH for go.
```
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

##2. install Google App Engine SDK for Go
download SDK and install.

set PATH for GAE/Go
```
export PATH=/path/to/go_appengine:$PATH
```

# files
## /urls.go
url rule

## /main.go
main page serve functions

## /models.go
Datastore models

## /admin.go
Admin applications, REST API for JS applications

## /settings.go
Application Settings

## /templates
html templates

## /media
static files
