package main

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"runtime"
)

func calcCrc32(data string, out chan string){
	out <- DataSignerCrc32(data)
}

func calcCrc32Md5(data string, out chan string, lock chan struct{}){
	lock <- struct{}{}
	dataMd5 := DataSignerMd5(data)
	<- lock
	out <- DataSignerCrc32(dataMd5)
}

func SingleHashWorker(data string, out chan interface{}, wg *sync.WaitGroup, lock chan struct{}) {
	defer wg.Done()

	crc32out := make(chan string, 1)
	go calcCrc32(data, crc32out)

	crc32Md5out:= make(chan string, 1)

	go calcCrc32Md5(data, crc32Md5out, lock)

	crc32Data := <- crc32out
	crc32Md5Data := <- crc32Md5out

	//fmt.Println("SingleHash crc32_md5_data:", crc32Md5Data)
	//fmt.Println("SingleHash crc32_data:", crc32Data)

	out <- crc32Data + "~" + crc32Md5Data
	runtime.Gosched()
}

func SingleHash(in, out chan interface{}) {
	fmt.Println("SingleHash")
	wg := &sync.WaitGroup{}
	var lock = make(chan struct{},1)
	for rawData := range in {
		int_data, ok := rawData.(int)
		data := strconv.Itoa(int_data)
		//fmt.Println("SingleHash data:", data)
		if !ok {
			panic("Wrong type in SingleHash")
		}
		wg.Add(1)
		go SingleHashWorker(data, out, wg, lock)
	}
	wg.Wait()
	close(out)
}

type multiHashResult struct{
	data string
	iter int
}

func calcThDataCrc32(th,data string, out chan multiHashResult){


	//out <- "test"
	iter, _ := strconv.Atoi(th)
	res := multiHashResult{
		data: DataSignerCrc32(th + data),
		iter: iter,
	}
	out <- res
}

type ByTh []multiHashResult

func (a ByTh) Len() int           { return len(a) }
func (a ByTh) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTh) Less(i, j int) bool { return a[i].iter < a[j].iter }


func MultiHashWorker(data string, out chan interface{},wg *sync.WaitGroup ){
	defer wg.Done()
	res := ""
	//step := ""
	var thString string
	//var results []multiHashResult
	results := make([]multiHashResult, 0, 50)
	workerOut := make(chan multiHashResult, 6)
	//fmt.Println("MultiHashWorker start: ",data)

	for th:=0; th<6; th++{
		thString = strconv.Itoa(th)
		go calcThDataCrc32(thString, data, workerOut)
		//step = DataSignerCrc32(thString + data)
		//fmt.Printf("MultiHash iter: %v step: %v \n",thString, step)
	}
	//close(workerOut)
	for th:=0; th<6; th++{
		step := <- workerOut
		results = append(results,step)
		//res = res + step.data
	}

	sort.Sort(ByTh(results))
	for _, st := range results{
		//step := <- workerOut
		res = res + st.data
	}
	//fmt.Println("MultiHashWorker done: ",res)
	out <- res
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for rawData := range in {

		data, ok := rawData.(string)
		//fmt.Println("MultiHash data:",data)
		if !ok {
			panic("Wrong type in SingleHash")
		}
		wg.Add(1)
		go MultiHashWorker(data, out, wg)
	}
	wg.Wait()
	close(out)
}

func CombineResults(in, out chan interface{}) {
	//fmt.Println("CombineResults")
	res := ""
	results := make([]string, 0, 50)
	for data := range in {
		//fmt.Println("CombineResults recevied", data)
		results = append(results, data.(string))
	}
	sort.Strings(results)
	for _, s := range results{
		res += s
		res += "_"
	}
	out <- res[:len(res)-1]
	close(out)
}

func ExecutePipeline(hashSignJobs ... job) {
	runtime.GOMAXPROCS(0)
	in := make(chan interface{}, 7)
	out := make(chan interface{}, 7)
	for ix, fn := range hashSignJobs {
		if ix < len(hashSignJobs) - 1 {
			go fn(in, out)
		} else {
			//fmt.Println("Final called")
			fn(in, out)
		}
		in = out
		out = make(chan interface{}, 7)
		//fmt.Printf("idx %v\n", ix)
	}
	fmt.Scanln()
}

func main() {
	println("run as\n\ngo test -v -race")
}
