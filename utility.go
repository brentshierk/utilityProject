package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Download struct{
Url string
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
		println("an error has occurred",err)
	}
	fmt.Println(resp.StatusCode)

	if resp.StatusCode > 299{
		println(resp.StatusCode)
		return errors.New(fmt.Sprintf("cant process response %v",resp.StatusCode))
	}

	size,err := strconv.Atoi(resp.Header.Get("Content-Length"))
	fmt.Printf("size is %v bytes \n",size)


	//fmt.Println( size)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Printf(string(size))
	//fmt.Printf("size is %v bytes",size)


	return nil
}

func (d Download) makeRequest(method string) (*http.Request, error)  {
	r,err := http.NewRequest(
		method,
		d.Url,
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
	fmt.Println(start)
	d := Download{
		Url : "https://raw.githubusercontent.com/brentshierk/Portfolio/master/src/router/index.js",
		targetPath: "zoom-test.zip",
		totalConnections: 10,
	}
	err := d.Do()
	if err != nil {
		print(err)
	}


}