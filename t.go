package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	//"encoding/base64"
	"bytes"
	"regexp"

	//"os"
	"io/ioutil"

	"github.com/PuerkitoBio/goquery"
)

var x int = 0
var currentTime = time.Now()

// writeImage encodes an image 'img' in jpeg format and writes it into ResponseWriter.
func writeImage(w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}

func blueHandler(w http.ResponseWriter, r *http.Request) {
	m := image.NewRGBA(image.Rect(0, 0, 240, 240))
	blue := color.RGBA{0, 0, 255, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

	var img image.Image = m
	writeImage(w, &img)

}

func httpGet(url string) (content string, statusCode int) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		statusCode = -100
		return
	}
	defer resp.Body.Close()
	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		statusCode = -200
		return
	}
	statusCode = resp.StatusCode
	content = string(data)
	return
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

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

func main() {

	http.HandleFunc("/blue/", blueHandler)

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(""))))

	http.HandleFunc("/hh", func(w http.ResponseWriter, r *http.Request) {

		//var resp , _ = httpGet("https://forum.gamer.com.tw/C.php?bsn=60076&snA=5040836")

		var urlStr string = "https://forum.gamer.com.tw/C.php?bsn=60076&snA=5041638"
		var rawCookies string = ""
		var resp string = doGet(urlStr, map[string]string{}, rawCookies)
		var dom, _ = goquery.NewDocumentFromReader(strings.NewReader(resp))

		var ar = dom.Find(".userid")
		var userid = ar.Eq(ar.Length() - 1).Text()

		/*dom.Find(".userid").Each(func(i int, selection *goquery.Selection) {
			//sn, _ := selection.Attr("id")
			userid = selection.Text()
			//fmt.Println(sn, "/", title)
		})*/

		//fmt.Println(resp);
		w.Header().Set(
			"Content-Type",
			"text/html",
		)

		io.WriteString(
			w,
			userid,
		)
	})

	http.HandleFunc("/img", func(w http.ResponseWriter, r *http.Request) {

		imagPath := "http://img2.bdstatic.com/img/image/166314e251f95cad1c8f496ad547d3e6709c93d5197.jpg"
		//圖片正則
		var reg, _ = regexp.Compile(`(\w|\d|_)*.jpg`)
		var name = reg.FindStringSubmatch(imagPath)[0]
		fmt.Print(name)
		//通過http請求獲取圖片的流文件
		var resp, _ = http.Get(imagPath)
		var body, _ = ioutil.ReadAll(resp.Body)
		//var out, _ = os.Create(name)
		//io.Copy(out, bytes.NewReader(body))

		var buffer = new(bytes.Buffer)

		io.Copy(buffer, bytes.NewReader(body))

		/*if err := jpeg.Encode(buffer, *img, nil); err != nil {
			log.Println("unable to encode image.")
		}*/

		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
		if _, err := w.Write(buffer.Bytes()); err != nil {
			log.Println("unable to write image.")
		}

	})

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {

		var mh01 = time.Now()
		var mh, _ = time.ParseDuration("3s")

		if currentTime.Add(mh).Before(mh01) {
			//處理邏輯
			fmt.Println("true")
			currentTime = time.Now()
		} else {
			fmt.Println("f")
		}

		//Add方法和Sub方法是相反的，獲取t0和t1的時間距離d是使用Sub，將t0加d獲取t1就是使用Add方法
		//k := time.Now()

		//一天之前
		//d, _ := time.ParseDuration("-1s")

		//一月之前
		fmt.Println(mh01.Sub(currentTime))

		x++
		var s string = "Hello World:" + strconv.Itoa(x) + " <br> " + (currentTime).Format("2006-01-02 15:04:05")

		//fmt.Fprintf(w,s);

		w.Header().Set(
			"Content-Type",
			"text/html",
		)

		io.WriteString(
			w,
			`<doctype html>
			<html>
				<head>
					<title>Hello World</title>
				</head>
				<body>
					`+s+`
				</body>
			</html>`,
		)

	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
