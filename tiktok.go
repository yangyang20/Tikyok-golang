package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/PuerkitoBio/goquery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
	pb "tiktok/proto"
	"tiktok/types"
	"time"
)

type Tiktok struct {
	videoIdList []string
	msToken     string
	Nickname    string
	videoList   []types.Video
	secUid      string
	downloadDir string
	cookies     []*http.Cookie
}

var client *http.Client
var config *types.Config

func init() {

	if _, err := toml.DecodeFile("./conf.toml", &config); err != nil {
		panic(err)
	}
	var uri *url.URL
	var gCurCookieJar *cookiejar.Jar
	uri, _ = url.Parse("http://127.0.0.1:1087")
	gCurCookieJar, _ = cookiejar.New(nil)
	client = &http.Client{
		//Timeout: time.Second * 3,
		Transport: &http.Transport{
			// 设置代理
			Proxy: http.ProxyURL(uri),
		},
		Jar: gCurCookieJar,
	}
}

func (t *Tiktok) httpRequest(urlS string, header map[string]string) *http.Response {
	req, _ := http.NewRequest("GET", urlS, nil)

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:100.0) Gecko/20100101 Firefox/100.0")
	if t.cookies == nil {
		req.Header.Set("Cookie", "msToken=TykS-sCqpMNLW4_svENXEmELp_ZWzEx50fdQJSqVzrZvMlgR-4IX3WgHmCbj9kCQOIKpzZADR0TJ4JaWVibS48Or7S3lhufjYaZ9agKxIFjI_O5hnr8vUKdYc_oWvLbTHfc11LLMSx1tPE4=; _abck=C7F012353D9E13A259FE981B6C7CC57E~-1~YAAQDH8auBOjBCGBAQAAN69mJAff4T12FSm8lnGCu3b7sVgSThLE0czUyGqKYO1mhLw384aCm1QwQ0PilypIBaPl1ZQkNvv9f6fs0EiLlr1lQ8InD6LDhR/8Gn+XwboPU38k2/70GxfVb/bVIs38g28eMCScPYvx83Q8c/w8Mp7uC9xn/C4HRcmDz9jUli5uWFNBCFa3favYsH9x8zXVagjRsC5ZCmF0lD2y/5etqcoQXxWKoRBwRcKlcgZcEZZjkrAvNE+Qy9p0/kWMyyp60lTUqYK1/jbMqhYIj399/O+4pcdXz5lw9EpkyWFx9+Auk1v5uROp2YQhbFGhffsN3CiHIYTjLwHU8Ixl/44/sUcDMC8NQD1pIsHc8IaR76q0NaID3pJpJ8Popg==~-1~-1~-1; ttwid=1%7C8VBu9S4iq2G7u3OM4e4sGvLhdJW5muQdw2r6tUH1h88%7C1654173136%7C200d7ff7d5279ea3bd964b867557783ede8a4643c4c7fbe858564da0ffd5f6c9; tt_csrf_token=YT2saYqC-oX8uR9LvkpMoHYi79H_DpW5wEN0; ak_bmsc=022DA513A416AA65EA1BDDFAA2F095DF~000000000000000000000000000000~YAAQDH8auBSjBCGBAQAAN69mJA+15wmxt5SRILnETnQL1QpklhoTUHt3NUuodCLP0XhCBKo0iRA79y0OssSMmUkOMgOvjx/Te7bUidcfDPPZ7+b54zw/aa3v1aJXS0WqhACPdq2uT6mR2uEQLggs9biGbSvbeG4bHsKir5YrYm8+hvupWLYRdlYRm0sThILNnLdAT41KkxjWIpCYn5OWK7goQiW044/g8C8obwQVA6EaVLh0gTxSL6FNzUWfi86N1l2LEBLSyTfBoW+FWLwC7RVJbj8YWHH+WCfiIdyShZXxryrk4mgaY7JtWzJsgXf1sDl4rjUT50sXKRKsPo4tUZfmNDA7ZrrrzX+jfgjMRwXLM42eRh0KDz85pEXv1gU3eYnKEJcYn74=; bm_sz=55B13B40A9F377E281D88FC6E376E9AC~YAAQDH8auBajBCGBAQAAN69mJA9iDkaqFfD8JDY1x5VwWnTs/F+32BdKgRCJoPK5Gftf3suZsfxmuAsRZn0CAdgXtluWj99RF3h7aPOTKmoPht13smEsTPUMnd2XrKUk/vUApiix0nhfiOKP9T5thVjcKi2PBLfBGLNyIJuShLhvey54tH+wy/9V2gbM9tuy+LKDs5MyN/TcGVOAdghOZgOXQ1IKK1dDxBrVnl6kDLsfD+y+R5fDRjkNx4mcHb6JQtbDXH/xt1zn9VMFJuw3v35E586LPMVlqbR7oxemc7Uvhm0=~4277296~3289667; bm_sv=7FBA0F0BFB4199D482EC738B52587955~YAAQDn8auPFF2/qAAQAAEn1rJA9AszoKCYrQ+sQS6cRu+9NVJi7gDGQe2SFhRm6gm3PJvKSS5Wy86t/SM7Iwl6TnnaiIcIBUDSGcYcaY2ZadPxiyBzX3Zo7kqmZd72WY5+mqwSd1rkCkspcqbTd2YZxhhIsrY4cCFSgHCnMPasfz96/5Kpu+luo9kKdVkZcFC5gMI45TBXHnjzMNMi8LP0HusxePSaxOMl2zjtm0n7wFblAy/CsaEjUxugVQBzKD~1; bm_mi=9C048D5A0700CF158ACCAE9263578A94~YAAQDH8auBijBCGBAQAAaOtmJA+o4OTTNDq85/elH100ymIUS69o0Le475L6bxtdh+u7sj8C/Z0Rd2CkSOAKU1cx4Bx8UImry2MS+0sAYl1apyPkMxfVMAfJYFO8qmthsAqkzI6SeHe6Uaqan2mJ6S3Af79QpM5zc2HGPxeMGVI9TjDe1mYfHV8tn4Df0MH8l5xhqXqkRzA68Th7vu1+crCnx6NqS1h4QYgDJ/z3TxgQgJZgk0RPF+Tow7gvqmhnU3tMCiTwCOqY+570TSR6jqFqtoEHSYo929TLnzpczdrwwrIgkLEqT1nhbwaBAHpAo51eBEjr1EgNdyMo7cyWMBLBYWaSxg==~1")
	} else {
		for _, cookie := range t.cookies {
			req.AddCookie(cookie)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		panic(urlS + "请求错误")
	}

	if t.cookies == nil {
		t.cookies = resp.Cookies()
		//return t.httpRequest(urlS, header)
	}
	//if t.msToken == "" {
	//	fmt.Println("req:", req)
	//
	//	for _, cookie := range resp.Cookies() {
	//		if cookie.Name == "msToken" && len(cookie.Value) == 148 {
	//			t.msToken = cookie.Value
	//			break
	//		}
	//	}
	//}
	return resp

}

func (t *Tiktok) setSecUid() {
	resp := t.httpRequest(config.Url, nil)
	if find := strings.Contains(config.Url, "www.tiktok.com"); find {
		doc, _ := goquery.NewDocumentFromReader(resp.Body)
		con := doc.Find("#SIGI_STATE").Text()
		content := new(types.HomePageContent)
		json.Unmarshal([]byte(con), &content)
		t.secUid = content.UserPage.SecUid
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		html := string(body)
		reg, _ := regexp.Compile(`[sec_uid|sec_user_id]=(.*?)&amp;`)
		regA := reg.FindStringSubmatch(html)
		t.secUid = regA[1]
	}
	defer resp.Body.Close()
}

func (t *Tiktok) setVideoIds(cursor string, ch chan int) {
	tt := map[string]string{
		"aid":              "1988",
		"app_language":     "zh-Hant-TW",
		"app_name":         "tiktok_web",
		"browser_language": "zh-CN",
		"browser_name":     "Mozilla",
		"browser_online":   "true",
		"browser_platform": "Linux x86_64",
		"browser_version":  "5.0 (X11)",
		"channel":          "tiktok_web",
		"cookie_enabled":   "true",
		"device_id":        "7103145171260917250",
		"device_platform":  "web_pc",
		"focus_state":      "true",
		"from_page":        "user",
		"history_len":      "2",
		"is_fullscreen":    "false",
		"is_page_visible":  "true",
		"os":               "linux",
		"priority_region":  "",
		"referer":          "",
		"region":           "TW",
		"screen_height":    "1080",
		"screen_width":     "1920",
		"tz_name":          "Asia/Shanghai",
		"webcast_language": "zh-Hant-TW",
		"count":            "30",
		"cursor":           cursor,
		"secUid":           t.secUid,
		//"secUid":   "MS4wLjABAAAA4JaVbcTcVqIVbmNWh9suQAwHD4T15deIwN42p8sxM9KnDE_c72Pnrr0i6_2IOQGl",
		"language": "zh-Hant-TW",
	}
	ttParams := map[string]string{
		"x-tt-params": t.grpcGenXTTParams(tt),
		//"x-tt-params": "0Ao2GQ2ts0LfwmhpM6WRLdtRH4xjmQiGrZl8+q7i9hLrJAiyAYWImGKX0LLma+1qImvm7CxK8DZKOvbkb6fL48rp112MJ+e2gqRai7TtTe2SiMfQLbpniuqY2/oBU87Y+5RgElYLZC4MIR7lAhfLrGOY373ljMVVa7Sw5IzYkP45GbWqP8iTTHjRaeJ2tBfFujwWXByjqAUbusrvAEh4lUI001Ejlx1QAYhUtvtPRfHIJUgAJDSNKbu6QumxhL3pcOk2xQBuEo41Ai3msQXQmPoXlu5BuOH//gRQVlIMOv3RteOJsE9hIE7mfAEXDbz5jnzdsDBhayJWdwuYmno1n9REXLux6KgNIjiQj77P47Vk8ZjaIc42B8W1deal7ihf8Jdyz7yBojvJQAAphBDBD4TtE1lWDd3C3cpskiCxLh1NsZoKareiGtORkwGV6OFNfFkpSXpUFNCLNG3jcHu9yN/IAeaGU9aHXn4hwNQu+ic152G+b72rbDf9YNyFls61FcFOcKhJtuMjAAPNFW7xK9SqscJDn67nShz/IAQGlwaF264S7RVCmPUgLZzrGu1/WsdRkNeoSjv6rnn6qIadeyekZu7foBNcguaD9KZOtfDVcIODl73wUv+vhGewRwXdVNrimaZPtwSAgClTyNnFgpsTAeXKCTtduh1Lui/YLbpe+60AAbP30Ekl0WaXBLvgVgV8OoobIzX2ZPJpBVRpXv3mZQtNDa0jDvG8euW+TbgVr/vJ3PFAySs0YwTYLFNPhnub83/QH9Fq3rRER6QcndsnfUBlahgaJ22PaRFQ6mc=",
	}
	params := "aid=1988&app_language=zh-Hant-TW&app_name=tiktok_web&" +
		"browser_language=zh-CN&browser_name=Mozilla&browser_online=true&browser_platform=Linux x86_64&" +
		"browser_version=5.0 (X11)&channel=tiktok_web&cookie_enabled=true&device_id=7103145171260917250&device_platform=web_pc&" +
		"focus_state=true&from_page=user&history_len=2&is_fullscreen=false&is_page_visible=true&os=linux&priority_region=&" +
		"referer=&region=TW&screen_height=1080&screen_width=1920&tz_name=Asia/Shanghai&webcast_language=zh-Hant-TW&" +
		"msToken=TykS-sCqpMNLW4_svENXEmELp_ZWzEx50fdQJSqVzrZvMlgR-4IX3WgHmCbj9kCQOIKpzZADR0TJ4JaWVibS48Or7S3lhufjYaZ9agKxIFjI_O5hnr8vUKdYc_oWvLbTHfc11LLMSx1tPE4=" +
		"&X-Bogus=DFSzKIVYvmXANJv6Sw7UjGIbA6yP&" +
		"_signature=_02B4Z6wo00f01fDz3VwAAIBDX9do5-fkIpHw49HAAB6n49"
	urlS := "https://t.tiktok.com/api/post/item_list/?" + url.QueryEscape(params)
	resp := t.httpRequest(urlS, ttParams)
	html, _ := ioutil.ReadAll(resp.Body)
	videoItemList := new(types.VideoItemList)
	err := json.Unmarshal(html, &videoItemList)
	if err != nil {
		fmt.Println("VideoIds 获取错误")
		return
	}
	if videoItemList.StatusCode == 0 {
		if videoItemList.HasMore {
			go t.setVideoIds(videoItemList.Cursor, ch)
		}
		for _, v := range videoItemList.ItemList {
			t.videoIdList = append(t.videoIdList, v.Id)
		}
		if !videoItemList.HasMore {
			ch <- 1
		}
	} else {
		fmt.Println(string(html))
		fmt.Println("videoId 获取失败")
		return
	}
	defer resp.Body.Close()
}

func (t *Tiktok) grpcGenXTTParams(paramsJson map[string]string) string {
	params, err := json.Marshal(paramsJson)
	if err != nil {
		fmt.Println("grpcGenXTTParams传递的参数有误")
		return ""
	}
	// Set up a connection to the server.
	conn, err := grpc.Dial(config.Grpc.Addr+":"+config.Grpc.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTikTokClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayEncryption(ctx, &pb.EncryptionRequest{TtParams: string(params)})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	//log.Printf("grpcServce Message: %s", r.GetTtParamsStr())
	return r.GetTtParamsStr()
}

func (t *Tiktok) setVideoListMain() {
	resIndex := 0 // 用来保存每次拆分的长度
	const size = 100.0
	arrNum := math.Ceil(float64(len(t.videoIdList)) / size)
	ch := make(chan int, int(arrNum))
	if len(t.videoIdList) < size {
		go t.setVideoList(strings.Join(t.videoIdList, ","), ch)
	} else {
		for i := 0; i < int(arrNum); i++ {
			go t.setVideoList(strings.Join(t.videoIdList[resIndex:size+resIndex], ","), ch)
		}
	}
	for i := 0; i < int(arrNum); i++ {
		<-ch
	}
}

func (t *Tiktok) setVideoList(videoIds string, ch chan int) {
	urlS := "https://api.tiktokv.com/aweme/v1/multi/aweme/detail/?aweme_ids=%5B" + videoIds + "%5D"
	fmt.Println("开始获取：", urlS)
	resp := t.httpRequest(urlS, nil)
	resB, _ := ioutil.ReadAll(resp.Body)
	videoDetail := new(types.VideoDetail)
	err := json.Unmarshal(resB, &videoDetail)
	if err != nil {
		fmt.Println("VideoDetail转换错误")
	}
	for _, detail := range videoDetail.AwemeDetails {
		downUrl := detail.Video.PlayAddr.UrlList[0]
		videoName := detail.Desc
		if t.Nickname == "" {
			t.Nickname = detail.Author.Nickname
		}
		t.videoList = append(t.videoList, types.Video{
			DownUrl: downUrl,
			Name:    videoName,
		})
	}
	defer resp.Body.Close()
	ch <- 1
}

func (t *Tiktok) download(v types.Video, ch chan int) {
	fmt.Println("开始下载：", v.Name)
	resp := t.httpRequest(v.DownUrl, nil)
	f, err := os.Create(t.downloadDir + "/" + v.Name + ".mp4")
	if err != nil {
		ch <- 0
		panic(err)
	}
	io.Copy(f, resp.Body)
	defer resp.Body.Close()
	defer fmt.Println("下载完成", v.Name)
	ch <- 1
}

func (t *Tiktok) DownloadVideo() {
	t.setSecUid()
	c := make(chan int, 1)
	t.setVideoIds("0", c)
	<-c
	if len(t.videoIdList) == 0 {
		fmt.Println("没有获取到videoIdList")
		return
	}
	t.setVideoListMain()
	t.downloadDir = config.DownloadPath + t.Nickname
	if _, err := os.Stat(t.downloadDir); os.IsNotExist(err) {
		os.MkdirAll(t.downloadDir, 0777)
	}
	l := len(t.videoList)
	if l == 0 {
		fmt.Println("没有获取到下载列表")
		return
	}
	ch := make(chan int, l)
	for i := 0; i < l; i++ {
		go t.download(t.videoList[i], ch)
	}
	for i := 0; i < l; i++ {
		fmt.Println(<-ch)
	}
}

func main() {
	t := new(Tiktok)
	t.DownloadVideo()

}
