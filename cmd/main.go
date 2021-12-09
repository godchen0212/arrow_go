package main

import (
	"example.com/arrow_go/storage"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"time"
)

func execute(i int, enableGC bool) {
	fmt.Println("round ", i)
	b, err := ioutil.ReadFile("parquet-arrow-example2.parquet") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(len(b))

	r, err2 := storage.NewPayloadReader(storage.DataType_Int64, b)
	if err2 != nil {
		panic(err2.Error())
	}
	fmt.Println(r)
	for j := 0; j < 10; j++ {
		int64s, _ := r.GetInt64FromPayload()
		fmt.Println(len(int64s))
	}
	err = r.ReleasePayloadReader()
	if err != nil {
		fmt.Println(err.Error())
	}
	if enableGC {
		runtime.GC()
	}
}

func main() {
	if len(os.Args) < 1 {
		_, _ = fmt.Fprint(os.Stderr, "usage: my_example --enableGC\n")
		return
	}

	os.Setenv("ARROW_DEFAULT_MEMORY_POOL", "mimalloc")

	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var enableGC bool
	flags.BoolVar(&enableGC, "enableGC", false, "enable manually gc")
	if err := flags.Parse(os.Args[1:]); err != nil {
		os.Exit(-1)
	}
	for i := 0; i < 100; i++ {
		execute(i, enableGC)
	}
	time.Sleep(10 * time.Second)
}
