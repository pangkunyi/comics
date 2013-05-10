/* vim: se ts=2 sw=2 enc=utf-8: */
package main

import(
	"kunyi.info/comics"
	"kunyi.info/weibo"
	"fmt"
	"bytes"
	"net/http"
)

func main(){
	http.HandleFunc("/", indexHandler)
	err :=http.ListenAndServe(":9002", nil)
	if err != nil{
		panic(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err :=r.ParseForm()
	if err != nil {
		panic(err)
	}
	signedCode := r.FormValue("signed_request")
	code, err :=weibo.ParseSignedRequest(signedCode)
	if err != nil{
		panic(err)
		fmt.Fprintf(w, "%p", err)
	}
	//url :=`https://api.weibo.com/oauth2/authorize?response_type=token&client_id=1710979917&redirect_uri=http%3A%2F%2Fapps.weibo.com%2Fkunyicomics`
	//http.Redirect(w, r, url, 303)
	fmt.Fprintf(w,"signed_request: %v\n %v", signedCode, code)
}

func postCalvinAndHobbes(){
	title:= "calvinandhobbes"
	status:= "#calvin&hobbes#"
	postComics(title, status)
}

func postComics(title, status string){
	gc := comics.GoComics{Title:title}
	gc.Init()
	picBytes, err := gc.PicBytes()
	if err != nil {
		panic(err)
	}

	var wb = &weibo.Weibo{}
	wb.Init()
	param := weibo.UploadParam{Pic:bytes.NewBuffer(picBytes), Status:status, Filename:title+".gif"}

	var resp weibo.UploadResponse
	_, err = wb.Upload(&param, &resp)
	if err != nil {
		panic (err)
	}
	fmt.Println(resp.Text)
}
