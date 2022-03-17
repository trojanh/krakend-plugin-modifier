package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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
	f(string(r)+"-response", r.responseDump, false, true)
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

type requestWrapper struct {
	method  string
	url     *url.URL
	query   url.Values
	path    string
	body    io.ReadCloser
	params  map[string]string
	headers map[string][]string
}

func (r requestWrapper) Method() string               { return r.method }
func (r requestWrapper) URL() *url.URL                { return r.url }
func (r requestWrapper) Query() url.Values            { return r.query }
func (r requestWrapper) Path() string                 { return r.path }
func (r requestWrapper) Body() io.ReadCloser          { return r.body }
func (r requestWrapper) Params() map[string]string    { return r.params }
func (r requestWrapper) Headers() map[string][]string { return r.headers }

type metadataWrapper struct {
	headers    map[string][]string
	statusCode int
}

func (m metadataWrapper) Headers() map[string][]string { return m.headers }
func (m metadataWrapper) StatusCode() int              { return m.statusCode }

type responseWrapper struct {
	data       map[string]interface{}
	isComplete bool
	metadata   metadataWrapper
	io         io.Reader
}

func (r responseWrapper) Data() map[string]interface{}            { return r.data }
func (r responseWrapper) IsComplete() bool                        { return r.isComplete }
func (r responseWrapper) Metadata() proxy.ResponseMetadataWrapper { return r.metadata }
func (r responseWrapper) Io() io.Reader                           { return r.io }

type InputUrl interface {
	Url() string
}

type InputWrapper interface {
	KrakendDebuggerResponse() InputUrl
}

func toString(body io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	return buf.String()
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
		fmt.Println("\n[REQUEST DUMP]: dumper injected!!!")

		req, ok := input.(RequestWrapper)
		if !ok {
			return nil, errors.New("request:something went wrong")
		}
		fmt.Println("\t[REQUEST DUMP]: url:", req.URL())
		for _, m := range cfg {
			fmt.Println("value is", m)
		}
		fmt.Println("\t[REQUEST DUMP]: params :", req.Params())
		// fmt.Println("\t[REQUEST DUMP]: headers:", req.Headers())
		// fmt.Println("\t[REQUEST DUMP]: method:", req.Method())
		fmt.Println("\t[REQUEST DUMP]: query:", req.Query())
		fmt.Println("\t[REQUEST DUMP]: body:", toString(req.Body()))
		// fmt.Println("\t[REQUEST DUMP]: path:", req.Path())

		return input, nil
	}
}

func callApi(url string) ( /*[]byte*/ string, error) {
	resp, err := http.Get(url)

	if err != nil {
		log.Printf("Request Failed: %s", err)
		// return []byte(""), err
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		// return []byte(""), err
		return "", err
	}
	// Log the request body
	bodyString := string(body)
	log.Print(bodyString)
	return bodyString, nil
}

func toReader(data string) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewBufferString(data))
}

func modifier(resp ResponseWrapper) responseWrapper {
	return responseWrapper{
		data:       resp.Data(),
		isComplete: resp.IsComplete(),
		metadata: metadataWrapper{
			headers:    resp.Metadata().Headers(),
			statusCode: resp.Metadata().StatusCode(),
		},
		io: resp.Io(),
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
		fmt.Println("\n[RESPONSE DUMP]: response dumper injected!!!")

		// tmp := cfg.(map[string]interface{})
		resp, ok := input.(ResponseWrapper)
		if !ok {
			return nil, errors.New("response:something went wrong")
		}

		// fmt.Println("data:", cfg["krakend-debugger-response"].(map[string]interface{})["url"].(string))
		// fmt.Println("is complete:", cfg)
		// strUrl, ok := cfg["name"]

		url := cfg["krakend-debugger-response"].(map[string]interface{})["url"].(string)
		// for key, m := range cfg {
		// 	fmt.Println("value is", m, key)
		// 	if key == "krakend-debugger-response" {
		// 		val, _ := m.(map[string]interface{})
		// 		url = fmt.Sprintln(val["url"])
		// 		fmt.Println("value is", val["url"])
		// 	}
		// }
		println("\t[RESPONSE DUMP]: CFG: ", url)

		respData := resp.Data()
		login := respData["org"].(map[string]interface{})["login"].(string)
		apiResp, err := callApi(url + login + ".json")

		// var mapResp map[string]interface{}
		// if err = json.Unmarshal(apiResp, &mapResp); err != nil {
		// 	return "", err
		// }
		if err != nil {
			log.Printf("Reading body failed: %s", err)
			return "", err
		}
		// input[".Data()"] = mapResp
		println("\t[RESPONSE DUMP]: API response:", apiResp)
		// fmt.Println("\t[RESPONSE DUMP]: data:", apiResp)
		fmt.Println("\t[RESPONSE DUMP]: data:", respData)
		fmt.Println("\t[RESPONSE DUMP]: headers:", resp.Metadata().Headers())
		fmt.Println("\t[RESPONSE DUMP]: status code:", resp.Metadata().StatusCode())
		// jsonStr := apiResp

		// Convert json string to map[string]interface{}
		var mapData map[string]interface{}
		if err := json.Unmarshal([]byte(apiResp), &mapData); err != nil {
			fmt.Println(err)
		}
		tmp := modifier(resp)

		tmp.data = mapData
		return tmp, nil
	}
}
