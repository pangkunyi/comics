/* vim: se ts=2 sw=2 enc=utf-8: */
package main

import(
	"kunyi.info/comics"
	"kunyi.info/weibo"
	"fmt"
	"bytes"
	"io"
	"net/http"
	"mime/multipart"
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
	param := weibo.UploadParam{Pic:&bytes.NewBuffer(picBytes), Status:status, Filename:"calvinandhobbes.gif"}

	var post weibo.WeiboPost
	_, err = weibo.Upload(req, &post)
	if err != nil {
		panic (err)
	}
	fmt.Println(post.Text)
}
