package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Download struct{
url string
targetPath string
totalConnections int 

}

func (d Download) Do() error{
	fmt.Println("making connection")
	r,err := d.makeRequest("HEAD")
	if err != nil {
		println(err)
	}
	resp,err := http.DefaultClient.Do(r)
	if err != nil {
		println(err)
	}
	fmt.Println(resp.StatusCode)

	if resp.StatusCode > 299{
		println(resp.StatusCode)
		return errors.New(fmt.Sprintf("cant process response %v",resp.StatusCode))
	}
	resp.Header.Get("Content-Length")
	return nil
}

func (d Download) makeRequest(method string) (*http.Request,error)  {
	r,err := http.NewRequest(
		method,
		d.url,
		nil,
		)
	if err != nil{
		println(err)
	}
	r.Header.Set("User-Agent","snag a file")
	return r,nil
}


func main(){
	fmt.Println("snag a file downloader")
	start := time.Now()
	d := Download{
		url : "https://drive.google.com/file/d/1W77cE-XINukX156f2AA7yHdEDCdavble/view"
		targetPath: "zoom.mp4",
		totalConnections: 10,
	}
	err := d.Do()


}