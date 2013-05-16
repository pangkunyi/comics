/* vim: se ts=2 sw=2 enc=utf8: */
package weibo

import (
	"fmt"
	"os"
	"io"
	"io/ioutil"
	"strings"
	"bytes"
	"errors"
	"net/http"
	"encoding/json"
)

const(
	PicUploadURL =`http://picupload.service.weibo.com/interface/pic_upload.php?app=miniblog&data=1&url=weibo.com/u/3428296462&markpos=1&logo=1&nick=%40%E7%9C%8B%E6%BC%AB%E7%94%BB-calvinandhobbes&marks=1&mime=image/gif&ct=0.4171791155822575`
	AddURL =`http://weibo.com/aj/mblog/add?_wv=5&__rnd=1368511377550`
	KeepAliveURL=`http://weibo.com/u/3428296462`
)

func httpCall(req *http.Request, v interface{}) (bool, error){
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, errors.New(fmt.Sprintf("req url error: %s, %v", req.RequestURI, err))
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body,err := ioutil.ReadAll(resp.Body)
		fmt.Println("req uri: ", req.URL.Path)
		if strings.Contains(req.URL.Path, "pic_upload"){
			bodyString := string(body)
			fmt.Println("resp body: ", bodyString)
			body = body[strings.LastIndex(bodyString, ">")+1:]
		}
		d := json.NewDecoder(bytes.NewBuffer(body))
		err = d.Decode(&v)
		if err != nil {
			return false, errors.New(fmt.Sprintf("encoding json error: %v",  err))
		}
		return true, nil
	}
	return false, err
}

func KeepAlive(){
	fmt.Printf("ready for keeping alive weibo....\n")
	req, err := http.NewRequest("GET", KeepAliveURL, nil)
	if err != nil {
		fmt.Printf("fail to make keep alive request, cause by: %v\n", err)
		return 
	}
	setHeader(req)
	req.Header.Set("Host","weibo.com")
	req.Header.Set("Origin","http://weibo.com")
	req.Header.Set("Referer","http://weibo.com/u/3428296462?wvr=5&topnav=1&wvr=5")
	req.Header.Set("X-Requested-With","XMLHttpRequest")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("fail to keep alive, cause by: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("keep alive weibo done.\n")
}

func Add(pic string, v interface{}) (bool, error) {
	data :="text=%23calvin%26hobbes%23&pic_id="+pic+"&rank=0&rankid=&_surl=&hottopicid=550&location=home&module=stissue&_t=0"
	req, err := http.NewRequest("POST", AddURL, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return false, err
	}
	setHeader(req)
	req.Header.Set("Host","weibo.com")
	req.Header.Set("Origin","http://weibo.com")
	req.Header.Set("Referer","http://weibo.com/u/3428296462?wvr=5&topnav=1&wvr=5")
	req.Header.Set("X-Requested-With","XMLHttpRequest")
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")
	return httpCall(req, v)
}

func PicUpload(pic io.Reader, v interface{}) (bool, error) {
	req, err := http.NewRequest("POST", PicUploadURL, pic)
	if err != nil {
		return false, err
	}
	setHeader(req)
	req.Header.Set("Content-Type","application/octet-stream")
	req.Header.Set("Host","picupload.service.weibo.com")
	req.Header.Set("Origin","http://js.t.sinajs.cn")
	req.Header.Set("Referer","http://js.t.sinajs.cn/t5/home/static/swf/MultiFilesUpload.swf?version=559f4bc1f6266504")
	return httpCall(req, v)
}

func setHeader(req *http.Request){
	req.Header.Set("Accept","*/*")
	req.Header.Set("Accept-Charset","UTF-8,*;q=0.5")
	req.Header.Set("Accept-Encoding","gzip,deflate,sdch")
	req.Header.Set("Accept-Language","en-US,en;q=0.8")
	req.Header.Set("Connection","keep-alive")
	req.Header.Set("Cookie",loadCookie())
	req.Header.Set("User-Agent","Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.31 (KHTML, like Gecko) Chrome/26.0.1410.64 Safari/537.31")
}

func loadCookie() string{
	data, err := ioutil.ReadFile(os.Getenv("HOME")+"/.comics/weibo-cookie")
	if err != nil {
		fmt.Printf("fail to load cookie, cause by: %v\n", err)
		return ""
	}
	return string(data)
}
