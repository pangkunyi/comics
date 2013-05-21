/* vim: se ts=2 sw=2 enc=utf8: */
package weibo

import (
	"code.google.com/p/mahonia"
	"fmt"
	"io"
	"net/http"
	"io/ioutil"
	"bytes"
)

const(
	PicUploadURL =`http://picupload.service.weibo.com/interface/pic_upload.php?app=miniblog&data=1&url=weibo.com/u/3428296462&markpos=1&logo=1&nick=%40%E7%9C%8B%E6%BC%AB%E7%94%BB-calvinandhobbes&marks=1&mime=image/gif&ct=0.4171791155822575`
	AddURL =`http://weibo.com/aj/mblog/add?_wv=5&__rnd=1368511377550`
	KeepAliveURL=`http://weibo.com/u/3428296462`
)

var gbk = mahonia.NewDecoder("GBK")

func (ac *Account) KeepAlive(){
	fmt.Printf("ready for keeping alive weibo....\n")
	req, err := http.NewRequest("GET", KeepAliveURL, nil)
	if err != nil {
		fmt.Printf("fail to make keep alive request, cause by: %v\n", err)
		return 
	}
	ac.setHeader(req)
	req.Header.Set("Host","weibo.com")
	req.Header.Set("Origin","http://weibo.com")
	req.Header.Set("Referer","http://weibo.com/u/3428296462?wvr=5&topnav=1&wvr=5")
	req.Header.Set("X-Requested-With","XMLHttpRequest")
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("fail to keep alive, cause by: %v\n", err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("resp body: ", gbk.ConvertString(string(body)))
	fmt.Printf("keep alive weibo done.\n")
}

func (ac *Account) Add(pic string, v interface{}) (bool, error) {
	data :="text=%23calvin%26hobbes%23&pic_id="+pic+"&rank=0&rankid=&_surl=&hottopicid=550&location=home&module=stissue&_t=0"
	req, err := http.NewRequest("POST", AddURL, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return false, err
	}
	ac.setHeader(req)
	req.Header.Set("Host","weibo.com")
	req.Header.Set("Origin","http://weibo.com")
	req.Header.Set("Referer","http://weibo.com/u/3428296462?wvr=5&topnav=1&wvr=5")
	req.Header.Set("X-Requested-With","XMLHttpRequest")
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")
	return httpCall(req, v)
}

func (ac *Account) PicUpload(pic io.Reader, v interface{}) (bool, error) {
	req, err := http.NewRequest("POST", PicUploadURL, pic)
	if err != nil {
		return false, err
	}
	ac.setHeader(req)
	req.Header.Set("Content-Type","application/octet-stream")
	req.Header.Set("Host","picupload.service.weibo.com")
	req.Header.Set("Origin","http://js.t.sinajs.cn")
	req.Header.Set("Referer","http://js.t.sinajs.cn/t5/home/static/swf/MultiFilesUpload.swf?version=559f4bc1f6266504")
	return httpCall(req, v)
}

func (ac *Account) setHeader(req *http.Request){
	req.Header.Set("Accept","*/*")
	req.Header.Set("Accept-Charset","UTF-8,*;q=0.5")
	req.Header.Set("Accept-Language","en-US,en;q=0.8")
	req.Header.Set("Connection","keep-alive")
	req.Header.Set("Cookie",ac.Cookie)
	req.Header.Set("User-Agent","Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.31 (KHTML, like Gecko) Chrome/26.0.1410.64 Safari/537.31")
}
