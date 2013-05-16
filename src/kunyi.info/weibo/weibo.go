/* vim: set ts=2 sw=2 enc=utf-8: */
package weibo

import (
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	API_BASE_URL = "https://api.weibo.com/2"
	API_UPLOAD = API_BASE_URL + "/statuses/upload.json"
)
type Weibo struct {
	AccessToken string
}

type UploadParam struct {
	Status string
	Pic io.Reader
	Filename string
}

type UploadResponse struct {
	Created_At              string
	Id                      int64
	Mid                     string
	Text                    string
	Source                  string
	Trucated                bool
	Thumbnail_Pic           string
	Bmiddle_Pic             string
	Original_Pic            string
}

type WeiboError struct {
	Err        string `json:"error"`
	Error_Code int64
	Request    string
}

func (e *WeiboError) String() string {
	return fmt.Sprintf("weibo error:[code:%d, desc:%s, request:%s]", e.Error_Code, e.Err, e.Request)
}

func Request(req *http.Request, v interface{}) (bool, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, errors.New(fmt.Sprintf("fetch url error: %s, %v", req.RequestURI, err))
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		d := json.NewDecoder(resp.Body)
		err = d.Decode(&v)
		if err != nil {
			return false, errors.New(fmt.Sprintf("encoding json error: %v",  err))
		}
		return true, nil
	} else {
		bytes, _ := ioutil.ReadAll(resp.Body)
		var e WeiboError
		err = json.Unmarshal(bytes, &e)
		if err != nil {
			return false, errors.New(fmt.Sprintf("encoding json error: %v",  err))
		}
		return false, errors.New(fmt.Sprint(e))
	}
	return false, err
}

func (wb *Weibo) Upload(param *UploadParam, v interface{}) (bool, error) {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	w.WriteField("access_token",wb.AccessToken)
	w.WriteField("status",param.Status)
	wr, err := w.CreateFormFile("pic", param.Filename)
	if err != nil {
		return false, err
	}
	_, err = io.Copy(wr, param.Pic)
	if err != nil {
		return false, err
	}
	w.Close()
	req, err := http.NewRequest("POST", API_UPLOAD, buf)
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	return Request(req, v)
}
