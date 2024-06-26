[![Stories in Ready](https://badge.waffle.io/go-chef/chef.png?label=ready&title=Ready)](https://waffle.io/go-chef/chef)
[![Build Status](https://app.wercker.com/status/9cfd4b53ea24e0894904067f283e4cf8/s "wercker status")](https://app.wercker.com/project/bykey/9cfd4b53ea24e0894904067f283e4cf8)
[![Coverage Status](https://coveralls.io/repos/go-chef/chef/badge.png?branch=master)](https://coveralls.io/r/go-chef/chef?branch=master)

# Chef Server API Client Library in Golang
This is a Library that you can use to write tools to interact with the chef server. 

## Install
 
    go get github.com/go-chef/chef

## Test
  
    go get -t github.com/go-chef/chef
    go test -v github.com/go-chef/chef
    test_chef_server: kitchen verify  # integration tests

## SSL

  If you run into an SSL verification problem when trying to connect to a ssl server with self signed certs set up your config object with `SkipSSL: true`

## Usage
This example is setting up a basic client that you can use to interact with all the service endpoints (clients, nodes, cookbooks, etc. At [@chefapi](https://docs.chef.io/api_chef_server.html))
More usage examples can be found in the [examples](examples) directory.

```go
package main

import (
	"fmt"
	"os"
	"github.com/go-chef/chef"
)

func main() {
	// read a client key
	key, err := os.ReadFile("key.pem")
	if err != nil {
		fmt.Println("Couldn't read key.pem:", err)
		os.Exit(1)
	}

	// build a client
	client, err := chef.NewClient(&chef.Config{
		Name: "foo",
		Key:  string(key),
		// goiardi is on port 4545 by default. chef-zero is 8889
		BaseURL: "http://localhost:4545",
	})
	if err != nil {
		fmt.Println("Issue setting up client:", err)
	}

	// List Cookbooks
	cookList, err := client.Cookbooks.List()
	if err != nil {
		fmt.Println("Issue listing cookbooks:", err)
	}

	// Print out the list
	fmt.Println(cookList)
}
```

## Chef API Error Status
To get the error status and error message returned from calls to the Chef API Server
you can use ChefError to unwind the ErrorResponse and access the original http error.
These methods are available to get specific information from the the error.

* Error() returns a formatted error message with the URL and status code.
* StatusCode() returns the original return status code.
* StatusMsg() returns the error message extracted from the error message body.
* StatusText() returns the returned error message body, usually JSON.
* StatusMethod() returns the name of the method used for the request.
* StatusURL() returns the URL object used for the request.

## CONTRIBUTING

If you feel like contributing, great! Just fork the repo, make your
improvements, and submit a pull request. Tests would, of course, be appreciated.
Adding tests where there are no tests currently would be even more appreciated.
At least, though, try and not break anything worse than it is. Test coverage has
improved, but is still an ongoing concern.

## AUTHORS

|               |                                                |
|:--------------|:-----------------------------------------------|
|Jesse Nelson   |[@spheromak](https://github.com/spheromak)
|AJ Christensen |[@fujin](https://github.com/fujin)
|Brad Beam      |[@bradbeam](https://github.com/bradbeam)
|Kraig Amador   |[@bigkraig](https://github.com/bigkraig)
|Mark Gibbons   |[@mark](https://github.com/markgibbons)

## COPYRIGHT

Copyright 2013-2014, Jesse Nelson

## LICENSE

Like many Chef ecosystem programs, go-chef/chef is licensed under the Apache 2.0
License. See the LICENSE file for details.

Chef is copyright (c) 2008-2014 Chef, Inc. and its various contributors.

Thanks go out to the fine folks of Opscode and the Chef community for all their
hard work.
