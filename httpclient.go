// httpclient is a simple JSON-RPC over HTTP client implmentation.
package httpclient

import (
	"fmt"
	"log"
	"net/http"
	"io"
	"bytes"
	"encoding/json"
	"math/rand"
	"strings"
)

// JsonrpcHttpClient sends a JSON-RPC request over HTTP POST and returns the servers
// response.
type JsonrpcHttpClient struct {
	Id uint64
	Http struct {
		User string
		Password string
		Host string
		Port string
		Ssl bool
	}
}

// EncodeRequest formats parameters for a JSON-RPC client request as a json string.
func (j *JsonrpcHttpClient) EncodeRequest(method string, params interface{}) ([]byte, error) {
	var request struct {
		// A String containing the name of the method to be invoked.
		Method string `json:"method"`
		// Object to pass as request parameter to the method.
		Params interface{} `json:"params"`
		// Used to match response to request.
		Id uint64 `json:"id"`
		// This tells the rpc server what version of jsonrpc we're using.
		Jsonrpc string `json:"jsonrpc"`
	}

	j.Id = uint64(rand.Int63())

	request.Method = method
	request.Params = params
	request.Jsonrpc = "2.0"
	request.Id = j.Id

	return json.Marshal(request)
}

// DecodeResponse decodes the response body of a client request into
// the interface reply.
func (j *JsonrpcHttpClient) DecodeResponse(body io.Reader) (result map[string]interface{}, err error) {
	var response struct {
		// The actual response from the server.
		Result *json.RawMessage `json:"result"`
		Error  interface{}      `json:"error"`
		// Value informing the client of an error calling RPC.
		// This value is echoed back to help match responses to requests.
		Id     uint64           `json:"id"`
	}
	var jsonresult []interface{}
	result = nil

	// Decoding the json formated body of the HTTP response into our struct
	if err = json.NewDecoder(body).Decode(&response); err != nil {
		return
	}
	if response.Error != nil {
		err = fmt.Errorf("%v", response.Error)
		return
	}
	if response.Result == nil {
		err = fmt.Errorf("server returned no result")
		return
	}
	err = json.Unmarshal(*response.Result, &jsonresult)
	if err != nil {
		return
	}
	if response.Id != j.Id {
		err = fmt.Errorf("Id is incorrect: expected=%d returned=%d", j.Id, response.Id)
		return
	}
	defer func() {
		if e := recover(); e != nil {
			result = nil
			err = e.(error)
		}
	}()
	result = jsonresult[0].(map[string]interface{})
	return
}

func (j *JsonrpcHttpClient) Execute(method string, params interface{}) (map[string]interface{}, error) {

	jsonrequest, err := j.EncodeRequest(method, params)

	body := bytes.NewBuffer(jsonrequest)

	url := []string{}
	if j.Http.Ssl {
		url = append(url,"https://")
	} else {
		url = append(url,"http://")
	}

	url = append(url,
		[]string{
			j.Http.User,
			":",
			j.Http.Password,
			"@",
			j.Http.Host,
			":",
			j.Http.Port,
			"/targetrpc",
		}...
	)
	httprequest, err := http.NewRequest("POST", strings.Join(url,""), body)
	httprequest.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(httprequest)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}
	if response.Status != "200 OK"{
		return nil, fmt.Errorf("Server reports status: %s",response.Status)
	}

	return j.DecodeResponse(response.Body)
}
