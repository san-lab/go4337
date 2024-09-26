package rpccalls

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"

	"encoding/json"

	"github.com/san-lab/go4337/state"
)

func POSTCall(url, key string, data []byte) (result []byte, err error) {
	client := http.Client{}
	url, err = CombineURLWithKey(url, key)
	if err != nil {
		return nil, fmt.Errorf("could not combine URL with key: %v", err)
	}
	buf := bytes.NewBuffer(data)
	state.Log("POSTing", url, "with", string(data))
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Body = io.NopCloser(buf)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not do request: %v", err)
	}
	defer resp.Body.Close()
	resbts, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %w", err)
	}
	state.Log("RAW Response:", string(resbts))
	return resbts, nil
}

func ApiFreeHandCall(url, key, methodTemplate string, params ...interface{}) (result []byte, err error) {
	sparams := make([]any, len(params))
	for i, p := range params {
		if p == nil {
			continue
		}
		sparams[i], err = json.Marshal(p)
		if err != nil {
			return nil, fmt.Errorf("could not marshal param: %v", err)
		}
	}
	calldat := fmt.Sprintf(methodTemplate, sparams...)
	return POSTCall(url, key, []byte(calldat))

}

// Returns unparsed "Result" if successful, otherwise nil and either *APIError or error
func ApiCall(url, key string, ar *APIRequest, result interface{}) ([]byte, error) {
	data, err := ar.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("could not marshal APIRequest: %w", err)
	}

	bt, err := POSTCall(url, key, data)
	if err != nil {
		return nil, fmt.Errorf("Error in POSTCall: %v", err)
	}
	aresp := &APIRPCResponse{}
	err = json.Unmarshal(bt, aresp)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response: %w", err)
	}
	if aresp.Error.Code != 0 {
		return nil, &aresp.Error
	}
	if result != nil {
		err = json.Unmarshal(aresp.Result, result)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal result: %w", err)
		}
	}
	return aresp.Result, nil
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

func (ae *APIError) Error() string {
	return fmt.Sprintf("APIError: %s (%d)", ae.Message, ae.Code)
}

var ApiURLTemplates = map[string]string{}

func CombineURLWithKey(url, key string) (kurl string, err error) {
	if len(key) == 0 {
		return url, nil
	}
	if strings.Contains(url, "{{.}}") {
		turl := new(template.Template).Lookup(url)
		if turl == nil {
			turl, err = template.New(url).Parse(url)
			if err != nil {
				return url, err
			}
		}
		wr := &strings.Builder{}
		err = turl.Execute(wr, key)
		return wr.String(), err
	}

	if url[len(url)-1] == '/' {
		return url + key, nil
	}
	return url + "/" + key, nil
}
