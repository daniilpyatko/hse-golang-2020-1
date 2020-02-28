package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

const iterMultiHash = 6

func main() {}

func StartJob(curJob job, in chan interface{}, out chan interface{}, curWg *sync.WaitGroup) {
	defer curWg.Done()
	defer close(out)
	curJob(in, out)
}

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{}, 100)
	out := make(chan interface{}, 100)
	wg := &sync.WaitGroup{}
	for i := len(jobs) - 1; i >= 0; i-- {
		wg.Add(1)
		go StartJob(jobs[i], in, out, wg)
		out = in
		in = make(chan interface{}, 100)
	}
	wg.Wait()
}

type pair struct {
	c chan string
	s string
}

var CountMd5Chan chan pair

// Обрабатываем все запросы к DataSignerMd5 в порядке очереди, а не параллельно
func CountMd5Queue() {
	for cur := range CountMd5Chan {
		res := DataSignerMd5(cur.s)
		cur.c <- res
	}
}

func SafeCountMd5(curS string) string {
	cur := pair{
		c: make(chan string, 1),
		s: curS,
	}
	CountMd5Chan <- cur
	res := <-cur.c
	return res
}

func CountSingleHash(s string, out chan interface{}) {
	m1 := SafeCountMd5(s)
	tmp := func(ch chan string, s string) {
		ch <- DataSignerCrc32(s)
	}
	res1 := make(chan string, 1)
	res2 := make(chan string, 1)
	go tmp(res1, s)
	go tmp(res2, m1)
	p1 := <-res1 + "~" + <-res2
	out <- p1
}

func StartSingleHash(s string, wg *sync.WaitGroup, out chan interface{}) {
	defer wg.Done()
	CountSingleHash(s, out)
}

func SingleHash(in, out chan interface{}) {
	CountMd5Chan = make(chan pair, 200)
	go CountMd5Queue()
	wg := &sync.WaitGroup{}
	for data := range in {
		wg.Add(1)
		go StartSingleHash(strconv.Itoa(data.(int)), wg, out)
	}
	wg.Wait()
	close(CountMd5Chan)
}

func CountMultiHash(s string, out chan interface{}) {
	p1 := ""
	ar := make([]string, 6)
	wg := &sync.WaitGroup{}
	wg.Add(6)
	tmp := func(ind int, s string) {
		ar[ind] = DataSignerCrc32(strconv.Itoa(ind) + s)
		defer wg.Done()
	}
	for i := 0; i < iterMultiHash; i++ {
		go tmp(i, s)
	}
	wg.Wait()
	for i := 0; i < iterMultiHash; i++ {
		p1 += ar[i]
	}
	out <- p1
}

func StartMultiHash(s string, wg *sync.WaitGroup, out chan interface{}) {
	defer wg.Done()
	CountMultiHash(s, out)
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for data := range in {
		wg.Add(1)
		go StartMultiHash(data.(string), wg, out)
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	var res []string
	for data := range in {
		res = append(res, data.(string))
	}
	sort.Strings(res)
	p1 := strings.Join(res, "_")
	out <- p1
}
