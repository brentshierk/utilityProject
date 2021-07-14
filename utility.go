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
	if err != nil {
		fmt.Println(err)
	}

	//creating a 2d array
	var connections = make([][2]int,d.totalConnections)
	//spliting the original file into smaller chunks for the workers to handle
	fileChunks := size/d.totalConnections
	fmt.Printf("each chunk is %v bytes\n",fileChunks)


	//algorithm to make sure each section is starting at a new file byte
	for i := range connections{
		if i ==0{
			//starting byte of first worker
			connections[i][0] = 0
		}else {
			connections[i][0] = connections[i-1][1] + 1
		}

		if i < d.totalConnections-1{
			//ending byte of each worker
			connections[i][1] = connections[i][0] + fileChunks
		}else{
			connections[i][1] = size -1
		}
	}





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