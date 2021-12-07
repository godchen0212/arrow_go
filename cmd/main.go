package main

import (
	"example.com/arrow_go/storage"
	"fmt"
	"io/ioutil"
	"runtime"
	"time"
)

func execute(i int) {
	fmt.Println("Round", i)
	b, err := ioutil.ReadFile("parquet-arrow-example2.parquet") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(len(b))

	r, err2 := storage.NewPayloadReader(storage.DataType_Int64, b)
	if err2 != nil {
		panic(err.Error())
	}
	fmt.Println(r)
	for j := 0; j < 10; j++ {
		int64s, _ := r.GetInt64FromPayload()
		fmt.Println(len(int64s))
		//time.Sleep(time.Second)
	}
	//fmt.Println("Start to release R")
	err = r.ReleasePayloadReader()
	if err != nil {
		fmt.Println(err.Error())
	}
	runtime.GC()
}

func main() {
	for i := 0; i < 1000; i++ {
		execute(i)
	}
	time.Sleep(10 * time.Second)
}
