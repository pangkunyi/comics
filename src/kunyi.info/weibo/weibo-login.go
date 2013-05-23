/* vim: se ts=2 sw=2 enc=utf8: */
package weibo

import(
	"code.google.com/p/mahonia"
	"os"
	"strconv"
	"fmt"
	"math/big"
	"strings"
	"errors"
	"net/http"
	"encoding/json"
	"encoding/hex"
	"crypto/rsa"
	"crypto/rand"
	"io/ioutil"
)

const (
	PRE_LOGIN_URL = `http://login.sina.com.cn/sso/prelogin.php?entry=sso&callback=sinaSSOController.preloginCallBack&su=%s&rsakt=mod&client=ssologin.js(v1.4.5)`
	SSO_LOGIN_URL = `http://login.sina.com.cn/sso/login.php?client=ssologin.js(v1.4.5)`
)
type Account struct{
	UserName string
	Password string
	loginUrl string
	Cookie string
}

type PreLoginData struct{
	Servertime int64
	Nonce string
	Pubkey string
	Rsakv string
}

func (ac *Account) Init(){
	data, err := ioutil.ReadFile(os.Getenv("HOME")+"/.comics/account")
	if err != nil {
		fmt.Printf("fail to load account, cause by: %v\n", err)
		panic(err)
	}
	dataStrings :=strings.Split(string(data),":")
	ac.UserName=dataStrings[0]
	ac.Password=dataStrings[1]
	fmt.Printf("inited Account: %v\n", ac)
}

func (pld *PreLoginData) GetPwd(password string) string{
	var n big.Int
	n.SetString(pld.Pubkey, 16)
	pk := rsa.PublicKey{&n,65537}
	msg := fmt.Sprintf("%d\t%s\n%s",pld.Servertime,pld.Nonce,password)
	fmt.Println("msg:", msg)
	code, err :=rsa.EncryptPKCS1v15(rand.Reader, &pk, []byte(msg))
	if err != nil{
		panic(err)
	}
	hexCode := hex.EncodeToString(code)
	fmt.Println("code: ", hexCode)
	return hexCode
}

func (ac *Account) Login() error{
	var pld PreLoginData
	err := pld.Load(ac.UserName)
	if err != nil{
		return err
	}
	client := &http.Client{}
	data := `entry=weibo&gateway=1&from=&savestate=7&useticket=1&pagerefer=http%3A%2F%2Fweibo.com%2Fa%2Fdownload&vsnf=1&su=`+ac.UserName+`&service=miniblog&servertime=`+strconv.FormatInt(pld.Servertime,10)+`&nonce=`+pld.Nonce+`&pwencode=rsa2&rsakv=`+pld.Rsakv+`&sp=`+pld.GetPwd(ac.Password)+`&encoding=UTF-8&prelt=415&url=http%3A%2F%2Fweibo.com%2Fajaxlogin.php%3Fframelogin%3D1%26callback%3Dparent.sinaSSOController.feedBackUrlCallBack&returntype=META`
	fmt.Printf("form data: %v\n", data)
	req, err:=http.NewRequest("POST", SSO_LOGIN_URL+"?"+data,nil)// bytes.NewBuffer([]byte(data)))
	if err != nil{
		return err
	}
	SetHeader(req)
	resp, err := client.Do(req)
	if err != nil{
		return err
	}
	fmt.Printf("resp: %v\n", resp)
	defer resp.Body.Close()
	fmt.Println("sso login cookie: ",resp.Header["Set-Cookie"])
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return err
	}
	enc := mahonia.NewDecoder("GBK")
	bodyString :=enc.ConvertString(string(body))
	sIdx :=strings.Index(bodyString,`replace("`)
	if sIdx > 0 {
		ac.loginUrl =bodyString[sIdx+9:]
		sIdx = strings.Index(ac.loginUrl, `")`)
		if sIdx > 0{
			ac.loginUrl = ac.loginUrl[:sIdx]
			fmt.Println("loginUrl:", ac.loginUrl)
			if strings.Contains(ac.loginUrl, "ticket"){
				return ac.saveCookie()
			}
		}
	}
	return errors.New(fmt.Sprintf("fail to login, loginUrl: %s", ac.loginUrl))
}
func notRedirect(req *http.Request, via []*http.Request) error{
	return nil
}
func (ac *Account) saveCookie() error{
	tr := &http.Transport{}
	req, err:=http.NewRequest("GET", ac.loginUrl,nil)
	if err != nil{
		return err
	}
	SetHeader(req)
	resp, err := tr.RoundTrip(req)
	if err != nil{
		return err
	}
	ac.Cookie=""
	for _, cookie := range resp.Header["Set-Cookie"] {
		ac.Cookie = ac.Cookie + cookie[:strings.Index(cookie,";")+1]+" "
	}
	fmt.Println("weibo cookie:", ac.Cookie)
	return nil
}


func (pld *PreLoginData) Load(username string) error{
	url := fmt.Sprintf(PRE_LOGIN_URL, username)
	resp, err :=http.Get(url)
	if err != nil{
		return errors.New(fmt.Sprintf("fail to fetch prelogin.php, cause by: %v", err))
	}
	body, err :=ioutil.ReadAll(resp.Body)
	if err != nil{
		return errors.New(fmt.Sprintf("fail to reading prelogin.php response body, cause by: %v", err))
	}
	defer resp.Body.Close()
	fmt.Println("PreLogin Cookie: ",resp.Header["Set-Cookie"])
	jsonBody := string(body)
	jsonBody = jsonBody[strings.Index(jsonBody,"(")+1:len(jsonBody)-1]
	fmt.Printf("pre login json body: %s\n", jsonBody)
	err =json.Unmarshal([]byte(jsonBody), &pld)
	if err != nil{
		return errors.New(fmt.Sprintf("fail to unmarshal pre login json body: %s, cause by: %v", jsonBody, err))
	}
	return nil
}

func SetHeader(req *http.Request){
	req.Header.Set("Accept","*/*")
	req.Header.Set("Accept-Charset","UTF-8,*;q=0.5")
	req.Header.Set("Accept-Language","en-US,en;q=0.8")
	req.Header.Set("Connection","keep-alive")
	req.Header.Set("User-Agent","Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.31 (KHTML, like Gecko) Chrome/26.0.1410.64 Safari/537.31")
}
