package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var t01_tim_上次請求時間 time.Time = time.Now()
var t01_byte_img []byte

var t02_tim_上次請求時間 time.Time = time.Now()
var ar_baha_user = make(map[string]*baha_user)

type baha_user struct {
	ip     string
	ip2    string
	id     string //帳號
	name   string //名稱
	gp     string
	level  string //等級
	career string //職業
	race   string //種族
}

func main() {

	//建立 server
	http.HandleFunc("/Reply/t02", func(w http.ResponseWriter, r *http.Request) {

		//回傳圖片
		w.Header().Set("Content-Type", "image/svg+xml")

		var ip = RemoteIp(r)        // 取得IP
		var sn = r.FormValue("snA") //取得get的參數

		var tim_now time.Time = time.Now()       //目前時間
		var tim_3s, _ = time.ParseDuration("3s") //3秒

		//如果間隔大於3秒
		if t02_tim_上次請求時間.Add(tim_3s).Before(tim_now) {
			t02_tim_上次請求時間 = time.Now() //更新最後請求時間
			fmt.Println("重新請求")

			//請求網址
			var urlStr string = "https://forum.gamer.com.tw/C.php?page=81000&bsn=60076&snA=" + sn
			var rawCookies string = ""
			var resp string = doGet(urlStr, map[string]string{}, rawCookies)
			var dom, _ = goquery.NewDocumentFromReader(strings.NewReader(resp))

			var ar = dom.Find(".c-section")
			func_add_ar_baha_user(ar)

			//如果有上一頁的話，就連同上一頁的內容都抓
			var pagebtnA = dom.Find(".BH-pagebtnA").Eq(0).Find("a")
			if pagebtnA.Length() >= 3 {
				href, _ := pagebtnA.Eq(pagebtnA.Length() - 2).Attr("href")
				href = "https://forum.gamer.com.tw/C.php" + href
				fmt.Println("上一頁：" + href)

				//抓取上一頁的內容
				var urlStr2 string = href
				var rawCookies2 string = ""
				var resp2 string = doGet(urlStr2, map[string]string{}, rawCookies2)
				var dom2, _ = goquery.NewDocumentFromReader(strings.NewReader(resp2))
				var ar2 = dom2.Find(".c-section")
				func_add_ar_baha_user(ar2)

			}

		}

		//取得使用者的IP前3個數字，用來跟巴哈的IP進行比對
		var ip_xxx = "0.0.0.0"
		var ar_ip = strings.Split(ip, ".")
		var img_base64 = ""
		if len(ar_ip) == 4 {
			ip_xxx = ar_ip[0] + "." + ar_ip[1] + "." + ar_ip[2]
		}

		var txt = "你根本沒回文吧"
		//fmt.Println(ip_xxx)
		if _, ok := ar_baha_user[ip_xxx]; ok {
			txt = "你是" + ar_baha_user[ip_xxx].race + ar_baha_user[ip_xxx].career + "沒錯吧"
			img_base64 = img01
		} else {
			img_base64 = img02
		}

		io.WriteString(
			w,
			`<?xml version="1.0" encoding="utf-8"?>
			<svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" x="0px" y="0px"
				 viewBox="0 0 400 275" style="enable-background:new 0 0 400 275;" xml:space="preserve">
			<style type="text/css">
				.st0{fill:#FFFFFF;}
				.st1{font-family:'AdobeMingStd-Light-B5pc-H';}
				.st2{font-size:28px;}
			</style>
			<g id="圖層_2">
				<rect width="400" height="275"/>
			</g>
			<g id="圖層_3">
				<text transform="matrix(1 0 0 1 74 261)" class="st0 st1 st2">`+txt+`</text>
			</g>
			<g id="圖層_1">				
					<image style="overflow:visible;" width="300" height="169" xlink:href="`+img_base64+`" transform="matrix(1.3333 0 0 1.3333 0 0)">
				</image>
			</g>
			</svg>
			`,
		)

		/*var sum = ""
		for key, value := range ar_baha_user {
			sum += key + " " + value.ip + "<br>" +
				"id:" + value.id + "<br>" +
				"career:" + value.career + "<br>" +
				"level:" + value.level + "<br>" +
				"race:" + value.race + "<br>" +
				"gp:" + value.gp + "<hr>"
		}*/
		/*io.WriteString(
			w,
			`<doctype html>
			<html>
				<head>
					<title>2</title>
					<meta charset="utf-8" />
				</head>
				<body>
					<img src= "`+img01+`">
					<h3>`+txt+" "+ip+sum+`</h3>
				</body>
			</html>`,
		)*/

	})

	//建立 server
	http.HandleFunc("/Reply/t01.png", func(w http.ResponseWriter, r *http.Request) {

		var tim_now time.Time = time.Now()       //目前時間
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
					<h1>你在期待什麼啦？</h1>
				</body>
			</html>`,
		)
	})

	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	//在指定的port上面進行啟用server
	http.ListenAndServe(":"+port, nil)
}

//-------------------

const (
	XForwardedFor = "X-Forwarded-For"
	XRealIP       = "X-Real-IP"
)

//
// RemoteIp 返回遠程客戶端的 IP，如 192.168.1.1
//
func RemoteIp(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get(XRealIP); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get(XForwardedFor); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}

//
// 解析網頁，儲存到 ar_baha_user 裡面
//
func func_add_ar_baha_user(ar *goquery.Selection) {

	ar.Each(func(i int, selection *goquery.Selection) {
		if selection.Find(".edittime").Length() > 0 {

			u_ip, _ := selection.Find(".edittime").Eq(0).Attr("data-hideip") //記錄於巴哈的IP
			u_id := selection.Find(".userid").Eq(0).Text()                   //帳號
			u_gp, _ := selection.Find(".usergp").Eq(0).Attr("title")         //gp
			u_level := selection.Find(".userlevel").Eq(0).Text()             //等級
			u_career, _ := selection.Find(".usercareer").Eq(0).Attr("title") //職業
			u_race, _ := selection.Find(".userrace").Eq(0).Attr("title")     //種族

			//fmt.Println(u_race)
			var c = new(baha_user)
			c.ip = strings.Replace(u_ip, ".xxx", "", -1)
			c.id = u_id
			c.level = u_level
			c.race = u_race
			c.career = u_career
			c.gp = u_gp

			ar_baha_user[c.ip] = c
		}
	})
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
