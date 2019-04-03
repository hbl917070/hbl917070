package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var t01_tim_上次請求時間 = time.Now()
var t01_byte_img []byte

func main() {

	//建立 server
	http.HandleFunc("/Reply/t01.png", func(w http.ResponseWriter, r *http.Request) {

		var tim_now = time.Now()                 //目前時間
		var tim_3s, _ = time.ParseDuration("3s") //3秒

		//如果間隔低於3秒
		if t01_tim_上次請求時間.Add(tim_3s).Before(tim_now) == false {
			if t01_byte_img != nil {

				//使用上次的圖片進行回傳
				w.Header().Set("Content-Type", "image/png")
				w.Header().Set("Content-Length", strconv.Itoa(len(t01_byte_img)))
				if _, err := w.Write(t01_byte_img); err != nil {
					log.Println("無法寫圖像")
				}
				return
			}
		}

		//請求網址
		var urlStr string = "https://forum.gamer.com.tw/C.php?page=81000&bsn=60076&snA=5037743"
		var rawCookies string = ""
		var resp string = doGet(urlStr, map[string]string{}, rawCookies)
		var dom, _ = goquery.NewDocumentFromReader(strings.NewReader(resp))

		//取得最後一個回文的帳號
		var ar = dom.Find(".userid")
		var userid = ar.Eq(ar.Length() - 1).Text()
		var user_img_url string = func_取得勇照網址(userid)

		//下載圖片
		t01_byte_img = func_download_img(user_img_url).Bytes()

		//回傳圖片
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(t01_byte_img)))
		if _, err := w.Write(t01_byte_img); err != nil {
			log.Println("無法寫圖像")
		}

		fmt.Println("重新請求")
		t01_tim_上次請求時間 = time.Now() //更新最後請求時間

	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set(
			"Content-Type",
			"text/html",
		)

		io.WriteString(
			w,
			`<doctype html>
			<html>
				<head>
					<title>你在期待什麼啦</title>
					<meta charset="utf-8" />
				</head>
				<body>
					<h1>你在期待什麼啦</h1>
				</body>
			</html>`,
		)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

//
//
//
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

//
// 請求網頁，並且回傳已解析的html物件
//
func doGet(urlStr string, queryData map[string]string, rawCookies string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlStr, nil)
	checkErr(err)
	q := req.URL.Query()
	for k, v := range queryData {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	var header = http.Header{}
	header.Add("Cookie", rawCookies)
	header.Add("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header = header
	resp, err := client.Do(req)
	checkErr(err)
	defer resp.Body.Close()
	ret, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	return string(ret)
}

//
// 取得勇照的網址
//
func func_取得勇照網址(s_userid string) string {

	//https://avatar2.bahamut.com.tw/avataruserpic/j/e/jeff60316377/jeff60316377.png
	s_userid = strings.ToLower(s_userid)
	var t1 string = string(s_userid[0])
	var t2 string = string(s_userid[1])
	var url_user_img string = "https://avatar2.bahamut.com.tw/avataruserpic/" + t1 + "/" + t2 + "/" + s_userid + "/" + s_userid + ".png"
	return url_user_img
}

//
// 下載圖片
//
func func_download_img(s_img_url string) *bytes.Buffer {

	//通過http請求獲取圖片的流文件
	var resp, _ = http.Get(s_img_url)
	var body, _ = ioutil.ReadAll(resp.Body)
	var buffer *bytes.Buffer = new(bytes.Buffer)

	io.Copy(buffer, bytes.NewReader(body))
	return buffer
}
