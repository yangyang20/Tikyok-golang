package main

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"regexp"
	"strconv"
	"strings"
	"tiktok/types"
)

type Douyin struct {
	secUid      string
	Nickname    string
	videoList   []types.Video
	downloadDir string
}

var clientD *http.Client
var configD *types.Config

func init() {
	if _, err := toml.DecodeFile("./conf.toml", &configD); err != nil {
		panic(err)
	}
	var gCurCookieJar *cookiejar.Jar
	gCurCookieJar, _ = cookiejar.New(nil)
	clientD = &http.Client{
		Jar: gCurCookieJar,
	}
}

func (d *Douyin) httpRequest(urlS string) *http.Response {
	defer func() {
		if err := recover(); err != nil {
			// 打印异常，关闭资源，退出此函数
			fmt.Println(err)
		}
	}()

	req, _ := http.NewRequest("GET", urlS, nil)

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:100.0) Gecko/20100101 Firefox/100.0")

	resp, err := clientD.Do(req)
	if err != nil {
		fmt.Println(urlS + "请求错误")
		panic(err)
	}
	if resp.StatusCode == http.StatusForbidden {
		return d.httpRequest(resp.Request.URL.String())
	}

	if resp.StatusCode != 200 {
		fmt.Println(urlS + "请求错误")
	}

	return resp
}

func (d *Douyin) setSecUid() {
	if find := strings.Contains(configD.Url, "www.douyin.com"); find {
		fmt.Println(configD.Url)
		reg := regexp.MustCompile(`user\/([\d\D]*)`)
		regA := reg.FindStringSubmatch(configD.Url)
		d.secUid = regA[1]
	} else {
		resp := d.httpRequest(configD.Url)
		reg := regexp.MustCompile(`user\/([\d\D]*)`)
		regA := reg.FindStringSubmatch(resp.Request.URL.Path)
		d.secUid = regA[1]
	}
}

func (d *Douyin) setVideo(max_cursor int, ch chan int) {
	urlS := "https://www.iesdouyin.com/web/api/v2/aweme/post/?sec_uid=" + d.secUid + "&count=30&max_cursor=" + strconv.Itoa(max_cursor) + "&aid=1128&_signature=PDHVOQAAXMfFyj02QEpGaDwx1S&dytk="
	fmt.Println(urlS)
	resp := d.httpRequest(urlS)
	body, _ := ioutil.ReadAll(resp.Body)

	videoDetail := new(types.VideoList)
	err := json.Unmarshal(body, &videoDetail)
	if err != nil {
		//html := string(body)
		fmt.Println(err)
		fmt.Println("video获取错误")
		return
	}
	for _, detail := range videoDetail.AwemeList {
		downUrl := detail.Video.PlayAddr.UrlList[0]
		uri := detail.Video.PlayAddr.Uri
		videoName := detail.Desc
		if d.Nickname == "" {
			d.Nickname = detail.Author.Nickname
		}
		d.videoList = append(d.videoList, types.Video{
			DownUrl: downUrl,
			Name:    videoName,
			Uri:     uri,
		})
	}
	if videoDetail.HasMore {
		d.setVideo(videoDetail.MaxCursor, ch)
	} else {
		ch <- 1
	}
	defer resp.Body.Close()
}

func (d *Douyin) download(v types.Video, ch chan int) {
	defer func() {
		if err := recover(); err != nil {
			// 打印异常，关闭资源，退出此函数
			fmt.Println(err)
			ch <- -1
		}
	}()

	fmt.Println("开始下载：", v.Name)
	var resp *http.Response
	if configD.IsUHD {
		urlS := "https://aweme.snssdk.com/aweme/v1/play/?video_id=" + v.Uri + "&radio=1080p&line=0"
		resp = d.httpRequest(urlS)
		if resp.StatusCode != http.StatusOK {
			resp = d.httpRequest(v.DownUrl)
		}
	} else {
		resp = d.httpRequest(v.DownUrl)
	}

	f, err := os.Create(d.downloadDir + "/" + v.Name + ".mp4")
	if err != nil {
		ch <- 0
		panic(err)
	}
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		ch <- 0
		fmt.Println(v.Name + "下载错误")
		return
	}
	defer resp.Body.Close()
	defer fmt.Println("下载完成", v.Name)

	ch <- 1
}

func (d *Douyin) DownloadVideo() {
	d.setSecUid()
	chV := make(chan int, 1)
	d.setVideo(0, chV)
	<-chV
	d.downloadDir = configD.DownloadPath + d.Nickname
	if _, err := os.Stat(d.downloadDir); os.IsNotExist(err) {
		os.MkdirAll(d.downloadDir, 0777)
	}

	l := len(d.videoList)
	ch := make(chan int, l)
	for i := 0; i < l; i++ {
		go d.download(d.videoList[i], ch)
	}
	for i := 0; i < l; i++ {
		<-ch
	}
}

func main() {
	douyin := new(Douyin)
	douyin.DownloadVideo()
}
