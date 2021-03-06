package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

)

type Download struct{
Url string
targetPath string
totalConnections int

}
func server() {
	r := gin.Default()
	r.GET(
		"/ping",
		func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	},
	)
	err := r.Run(":3000")
	if err != nil {
		return
	}
}

// Do func, makes a HEAD request to host server, the response is is the header where we can find the content-length(bytes)
//with the bytes we are then able to
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
	//splitting the original file into smaller chunks to make use of bandwidth thus potentially giving us faster download times
	// as the process is split amongst multiple threads
	fileChunks := size/d.totalConnections
	fmt.Printf("each chunk is %v bytes\n",fileChunks)


	//algorithm to make sure each section is starting at a new byte
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
	log.Println(connections)
	//using concurrency to download each section of the file
	// implementation of waitgroup https://tutorialedge.net/golang/go-waitgroup-tutorial/
	var waitgroup sync.WaitGroup
	for i,c := range connections{
		waitgroup.Add(1)
		//starting go routine
		go func(i int,c [2]int) {
			defer waitgroup.Done()
			err = d.downloadChunks(i,c)
		}(i,c)
	}
	waitgroup.Wait()
	err = d.mergeFileChunks(connections)
	if err != nil {
		return err
	}
	return d.mergeFileChunks(connections)
}
//helper function
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
//download the file chunks into separate tmp files
func (d Download) downloadChunks(i int, c [2]int) error {
	r,err := d.makeRequest("GET")
	if err != nil {
		return err
	}
	r.Header.Set("Range",fmt.Sprintf("bytes=%v-%v",c[0],c[1]))
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	if resp.StatusCode >299{
		return errors.New(fmt.Sprintf("error! response is %v",resp.StatusCode))
	}
	fmt.Printf("downloaded %v bytes for section %v \n",resp.Header.Get("Content-Length"),i)
	b,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf("fileChunk-%v.tmp",i),b,os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
//this function merges all of the separate file chunks into a single file.
func (d Download) mergeFileChunks(fileChunks [][2]int) error {
	f,err := os.OpenFile(d.targetPath,os.O_CREATE|os.O_WRONLY|os.O_APPEND,os.ModePerm)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	//iterate over the filechunks and write to file
	for i := range fileChunks{
		b,err := ioutil.ReadFile(fmt.Sprintf("fileChunk-%v.tmp",i))
		if err != nil {
			return err
		}
		n,err := f.Write(b)
		if err != nil {
			return err
		}
		fmt.Printf("%v bytes merged\n",n)
	}
	return nil
}


func main(){

	fmt.Println("snag a file downloader")
	start := time.Now()
	fmt.Println(start)
	d := Download{
		Url : "https://raw.githubusercontent.com/brentshierk/Portfolio/master/src/router/index.js",
		targetPath: "index.js",
		totalConnections: 10,
	}
	err := d.Do()
	if err != nil {
		print(err)
	}
	server()

}