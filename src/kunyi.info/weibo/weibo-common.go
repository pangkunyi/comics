/* vim: set ts=2 sw=2 enc=utf-8: */
package weibo

import(
	"fmt"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"encoding/json"
	"os"
	"io/ioutil"
	"strings"
)

var (
	clientSecret=loadClientSecret()
)

func base64Decode(payload string) ([]byte, error){
	for i:=0;i<4;i++{
		_bytes, err := base64.URLEncoding.DecodeString(payload)
		if err != nil {
			if i>3 {
				return nil, err
			}
		}else{
			return _bytes, nil
		}
		payload = payload+"="
	}
	return nil,errors.New("unknown error")
}

func ParseSignedRequest(input string) (map[string]interface{}, error){
	entries := strings.SplitN(input, ".", 2)
	entries[0]=entries[0]
	entries[1]=entries[1]
	buf, err := base64Decode(entries[0])
	if err != nil {
		return nil, err
	}
	sig := string(buf)

	buf, err = base64Decode(entries[1])
	if err != nil {
		return nil, err
	}
	var data interface{}
	err = json.Unmarshal(buf, &data)
	if err != nil {
		return nil, err
	}
	dataM, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New(fmt.Sprintf("unmarshal json error: %v", string(buf)))
	}
	if dataM["algorithm"] != "HMAC-SHA256" {
		return nil, errors.New(fmt.Sprintf("algorithm error: %v", dataM["algorithm"]))
	}
	mac :=hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(entries[1]))
	digest :=string(mac.Sum(nil))
	if digest != sig {
		return nil, errors.New(fmt.Sprintf("digest error: actual: %v, expected: %v", digest, sig))
	}
	return dataM, nil
}

func loadClientSecret() string{
	data, err := ioutil.ReadFile(os.Getenv("HOME")+"/.comics/weibo-client-secret")
	if err != nil {
		fmt.Printf("fail to load cookie, cause by: %v\n", err)
		return ""
	}
	return string(data)
}
