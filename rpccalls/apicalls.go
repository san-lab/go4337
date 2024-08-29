package rpccalls

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"encoding/json"
)

func POSTCall(url, key string, data []byte) (result []byte, aerr *APIError, err error) {
	client := http.Client{}
	url = url + "/" + key
	buf := bytes.NewBuffer(data)
	//fmt.Println("POSTing", url, "with", string(data))
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Body = io.NopCloser(buf)
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("could not do request: %v", err)
	}
	defer resp.Body.Close()
	resbts, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("could not read response: %v", err)
	}
	aresp := &APIRPCResponse{}
	err = json.Unmarshal(resbts, aresp)
	if err != nil {
		return nil, nil, fmt.Errorf("could not unmarshal response: %v", err, string(resbts))
	}
	if aresp.Error.Code != 0 {
		return nil, &aresp.Error, nil
	}
	return aresp.Result, nil, nil
}

func ApiFreeHandCall(url, key, methodTemplate string, params ...interface{}) (result []byte, aerr *APIError, err error) {
	sparams := make([]any, len(params))
	for i, p := range params {
		sparams[i], err = json.Marshal(p)
		if err != nil {
			return nil, nil, fmt.Errorf("could not marshal param: %v", err)
		}
	}
	calldat := fmt.Sprintf(methodTemplate, sparams...)
	return POSTCall(url, key, []byte(calldat))

}

// Returns unparsed "Result" if successful, otherwise nil and either *APIError or error
func ApiCall(url, key string, ar *APIRequest) ([]byte, *APIError, error) {
	data, err := ar.ToJSON()
	if err != nil {
		return nil, nil, fmt.Errorf("could not marshal APIRequest: %v", err, ar)
	}

	return POSTCall(url, key, data)

}

type APIRequest struct {
	ID      int           `json:"id"`
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func (ar *APIRequest) ToJSON() ([]byte, error) {
	return json.Marshal(ar)
}

type APIRPCResponse struct {
	ID      int             `json:"id"`
	Jsonrpc string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   APIError        `json:"error,omitempty"`
}

type APIError struct {
	Code    int             `json:"code"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
}
