package main

import (
	"example.com/arrow_go/storage"
	"fmt"
	//"io/ioutil"
	"os"
	//"runtime"
	"time"
)

func execute(i int, enableGC bool) {
	fmt.Println("round ", i)
	/*
		b, err := ioutil.ReadFile("parquet-arrow-example.parquet") // just pass the file name
		if err != nil {
			fmt.Print(err)
		}
		fmt.Println(len(b))
	*/

	_, err2 := storage.NewPayloadReader(storage.DataType_Int64, nil)
	if err2 != nil {
		panic(err2.Error())
	}
	/*
		fmt.Println(r)
			for j := 0; j < 10; j++ {
				int64s, _ := r.GetInt64FromPayload()
				fmt.Println(len(int64s))
			}
			err = r.ReleasePayloadReader()
			if err != nil {
				fmt.Println(err.Error())
			}
	*/
}

func main() {
	if len(os.Args) < 1 {
		_, _ = fmt.Fprint(os.Stderr, "usage: my_example --enableGC\n")
		return
	}

	for i := 0; i < 100; i++ {
		execute(i, false)
	}
	//runtime.GC()
	time.Sleep(1000 * time.Second)
}
