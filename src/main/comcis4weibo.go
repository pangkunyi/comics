/* vim: se ts=2 sw=2 enc=utf-8: */
package main

import(
	"kunyi.info/comics"
	"kunyi.info/weibo"
	"fmt"
	"bytes"
)

func main(){
	gc := comics.GoComics{Title:"calvinandhobbes"}
	gc.Init()
	picBytes, err := gc.PicBytes()
	if err != nil {
		panic(err)
	}

	var wb = &weibo.Weibo{}
	wb.Init()
	status :="#calvin&hobbes#"
	param := weibo.UploadParam{Pic:bytes.NewBuffer(picBytes), Status:status, Filename:"calvinandhobbes.gif"}

	var resp weibo.UploadResponse
	_, err = wb.Upload(&param, &resp)
	if err != nil {
		panic (err)
	}
	fmt.Println(resp.Text)
}
