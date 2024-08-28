package rpccalls

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"encoding/json"

	"github.com/san-lab/go4337/userop"
)

const StdBundlerCallTemplate = `curl -X GET %s/%s    --header 'accept: application/json'      --header 'content-type: application/json'      -d '
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "eth_supportedEntryPoints",
  "params": []
}
'
`

func ApiCall(url, key string, methodTemplate string, usop *userop.UserOperation) ([]byte, error) {
	client := http.Client{}
	url = url + "/" + key
	ausop := usop.ToUserOpForApi()
	data, err := json.Marshal(ausop)
	if err != nil {
		return nil, fmt.Errorf("could not marshal userop: %v", err)
	}

	methodCall := fmt.Sprintf(methodTemplate, string(data))
	//fmt.Println("methodCall:", methodCall)
	buf := strings.NewReader(methodCall)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		//fmt.Println(err)
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Body = io.NopCloser(buf)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not do request: %v", err)
	}
	return io.ReadAll(resp.Body)

}
