package main

import (
	"errors"
	"fmt"
	"io"
	"net/url"

	"github.com/luraproject/lura/v2/proxy"
)

func main() {}

func init() {
	fmt.Println(string(ModifierRegisterer), "loaded!!!")
}

// ModifierRegisterer is the symbol the plugin loader will be looking for. It must
// implement the plugin.Registerer interface
// https://github.com/luraproject/lura/blob/master/proxy/plugin/modifier.go#L71
var ModifierRegisterer = registerer("krakend-debugger")

type registerer string

// RegisterModifiers is the function the plugin loader will call to register the
// modifier(s) contained in the plugin using the function passed as argument.
// f will register the factoryFunc under the name and mark it as a request
// and/or response modifier.
func (r registerer) RegisterModifiers(f func(
	name string,
	factoryFunc func(map[string]interface{}) func(interface{}) (interface{}, error),
	appliesToRequest bool,
	appliesToResponse bool,
)) {
	f(string(r)+"-request", r.requestDump, true, false)
	f(string(r), r.responseDump, false, true)
	fmt.Println(string(r), "registered!!!")
}

// RequestWrapper is an interface for passing proxy request between the krakend pipe
// and the loaded plugins
type RequestWrapper interface {
	Params() map[string]string
	Headers() map[string][]string
	Body() io.ReadCloser
	Method() string
	URL() *url.URL
	Query() url.Values
	Path() string
}

// ResponseWrapper is an interface for passing proxy response between the krakend pipe
// and the loaded plugins
type ResponseWrapper interface {
	Data() map[string]interface{}
	Io() io.Reader
	IsComplete() bool
	Metadata() proxy.ResponseMetadataWrapper
}

func (r registerer) requestDump(
	cfg map[string]interface{},
) func(interface{}) (interface{}, error) {
	// check the cfg. If the modifier requires some configuration,
	// it should be under the name of the plugin.
	// ex: if this modifier required some A and B config params
	/*
	   "extra_config":{
	       "plugin/req-resp-modifier":{
	           "name":["krakend-debugger"],
	           "krakend-debugger":{
	               "A":"foo",
	               "B":42
	           }
	       }
	   }
	*/

	// return the modifier
	fmt.Println("request dumper injected!!!")
	return func(input interface{}) (interface{}, error) {
		fmt.Println("request dumper injected!!!")

		req, ok := input.(RequestWrapper)
		if !ok {
			return nil, errors.New("request:something went wrong")
		}

		fmt.Println("params:", req.Params())
		fmt.Println("headers:", req.Headers())
		fmt.Println("method:", req.Method())
		fmt.Println("url:", req.URL())
		fmt.Println("query:", req.Query())
		fmt.Println("path:", req.Path())

		return input, nil
	}
}

func (r registerer) responseDump(
	cfg map[string]interface{},
) func(interface{}) (interface{}, error) {
	// check the cfg. If the modifier requires some configuration,
	// it should be under the name of the plugin.
	// ex: if this modifier required some A and B config params
	/*
	   "extra_config":{
	       "plugin/req-resp-modifier":{
	           "name":["krakend-debugger"],
	           "krakend-debugger":{
	               "A":"foo",
	               "B":42
	           }
	       }
	   }
	*/

	// return the modifier
	fmt.Println("response dumper injected!!!")
	return func(input interface{}) (interface{}, error) {
		fmt.Println("response dumper injected!!!")

		resp, ok := input.(ResponseWrapper)
		if !ok {
			return nil, errors.New("response:something went wrong")
		}

		fmt.Println("data:", resp.Data())
		fmt.Println("is complete:", resp.IsComplete())
		fmt.Println("headers:", resp.Metadata().Headers())
		fmt.Println("status code:", resp.Metadata().StatusCode())

		return input, nil
	}
}
