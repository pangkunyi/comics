/* vim: se ts=1 sw=2 enc=utf-8: */
package main

import(
	"kunyi.info/comics"
	"io/ioutil"
	"os"
	"kunyi.info/weibo"
	"fmt"
	"bytes"
	"time"
)

var ac weibo.Account
func main(){
	ac.Init()
	for {
		postCalvinAndHobbes()
		fmt.Println("sleep 5 min for check comics update")
		time.Sleep(5 * 60 * time.Second)
	}
}

func postCalvinAndHobbes(){
	title:= "calvinandhobbes"
	status:= "#calvin&hobbes#"
	postComics(title, status)
}

func postComics(title, status string){
	gc := comics.GoComics{Title:title}
	gc.Init()
	err :=gc.Parse()
	if err != nil{
		fmt.Printf("can not parse pic url from gocomics")
		return
	}
	if gc.PicUrl == loadLastPicUrl() {
		fmt.Println("comics not update")
		return
	}
	err = ac.Login()
 	if err != nil{
 		fmt.Printf("weibo login error: %v\n",err)
		return
 	}
	picBytes, err := gc.PicBytes()
	if err != nil {
		fmt.Printf("fail to get pic data, cause of: %v\n", err)
		return
	}

	var resp interface{}
	_, err = ac.PicUpload(bytes.NewBuffer(picBytes), &resp)
	if err != nil {
		fmt.Printf("fail to upload pic data to weibo, cause of: %v\n", err)
		return
	}
	picM :=resp.(map[string]interface{})

	if picM["code"].(string) == "A00006" {
		picId :=picM["data"].(map[string]interface{})["pics"].(map[string]interface{})["pic_1"].(map[string]interface{})["pid"].(string);
		_, err =ac.Add(picId, &resp)
		if err != nil {
			fmt.Printf("post to weibo, cause of: %v\n", err)
			return
		}
		if resp.(map[string]interface{})["code"].(string) != "100000" {
			fmt.Printf("post to weibo error: %v\n", resp)
			return
		}
	}else{
			fmt.Printf("fail to upload pic data to weibo, cause of: %v\n", resp)
			return
	}
	writeLastPicUrl(gc.PicUrl)
	fmt.Println("Done.")
}

func loadLastPicUrl() string{
	data, err := ioutil.ReadFile(os.Getenv("HOME")+"/.comics/last_pic_url")
	if err != nil {
		fmt.Printf("fail to load last pic url, cause by: %v\n", err)
		return ""
	}
	return string(data)
}

func writeLastPicUrl(lastPicUrl string) {
	err := ioutil.WriteFile(os.Getenv("HOME")+"/.comics/last_pic_url", []byte(lastPicUrl), os.ModePerm)
	if err != nil {
		fmt.Printf("fail to write last pic url, cause by: %v\n", err)
		return
	}
}
func writeCookie(cookie string) {
	err := ioutil.WriteFile(os.Getenv("HOME")+"/.comics/weibo-cookie", []byte(cookie), os.ModePerm)
	if err != nil {
		fmt.Printf("fail to write weibo cookie, cause by: %v\n", err)
		return
	}
}
