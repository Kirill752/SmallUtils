package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func parallelSort(arr []string, numWorkers int) []string {
	if len(arr) <= 1 {
		return arr
	}
	// Разделяем слайс на части
	chunkSize := len(arr) / numWorkers
	if chunkSize == 0 {
		chunkSize = 1
	}
	wg := &sync.WaitGroup{}
	chunks := make([][]string, numWorkers)

	// Запускаем параллельную сортировку
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			start := i * chunkSize
			end := start + chunkSize
			if i == numWorkers-1 {
				end = len(arr) - 1
			}
			if end > len(arr)-1 {
				end = len(arr) - 1
			}
			// Копируем часть слайса
			chunks[i] = make([]string, end-start)
			copy(chunks[i], arr[start:end])
			// Сортируем эту часть
			sort.Slice(chunks[i], func(m, k int) bool { return chunks[i][m] < chunks[i][k] })
		}(i)
	}
	wg.Wait()
	return mergeSortedChunks(chunks)
}

func mergeSortedChunks(chunks [][]string) []string {
	res := make([]string, 0)
	for _, chunk := range chunks {
		res = merge(res, chunk)
	}
	return res
}

func merge(a, b []string) []string {
	result := make([]string, 0, len(a)+len(b))
	i, j := 0, 0

	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			result = append(result, a[i])
			i++
		} else {
			result = append(result, b[j])
			j++
		}
	}

	// Добавляем оставшиеся элементы
	result = append(result, a[i:]...)
	result = append(result, b[j:]...)

	return result
}

func SingleHash(in, out chan interface{}) {
	// wg := &sync.WaitGroup{}
	for data := range in {
		// wg.Add(1)
		resDSC32 := make(chan string, 1)
		resMd5 := make(chan string, 1)
		in_DSMd5Chan := make(chan string, 1)
		DSMd5_DSC32Chan := make(chan string, 1)
		in_DSC32Chan := make(chan string, 1)
		// Обработчик DataSignerCrc32 data из in
		go func() {
			data := <-in_DSC32Chan
			res := DataSignerCrc32(data)
			resDSC32 <- res
			close(resDSC32)
		}()
		// Обработчик DataSignerCrc32 data из DSMd5
		go func() {
			data := <-DSMd5_DSC32Chan
			res := DataSignerCrc32(data)
			resMd5 <- res
			close(resMd5)
		}()
		// Обработчик DSMd5 data из in
		go func() {
			data := <-in_DSMd5Chan
			res := DataSignerMd5(data)
			DSMd5_DSC32Chan <- res
			close(DSMd5_DSC32Chan)
		}()

		res := strings.Builder{}
		val := strconv.Itoa(data.(int))
		in_DSC32Chan <- val
		in_DSMd5Chan <- val
		res.WriteString(<-resDSC32)
		res.WriteString("~")
		res.WriteString(<-resMd5)
		out <- res.String()
	}
}

func MultiHash(in, out chan interface{}) {
	wgM := &sync.WaitGroup{}
	for data := range in {
		wgM.Add(1)
		go func(data interface{}) {
			// quotaCh <- struct{}{}
			defer wgM.Done()
			wg := &sync.WaitGroup{}
			res := make([]string, 6)
			for i := 0; i < 6; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					res[i] = DataSignerCrc32(strconv.Itoa(i) + data.(string))
				}(i)
			}
			wg.Wait()
			out <- strings.Join(res, "")
			// <-quotaCh
		}(data)
	}
	wgM.Wait()
}

func CombineResults(in, out chan interface{}) {
	res := []string{}
	for data := range in {
		res = append(res, data.(string))
	}
	sort.Slice(res, func(i, j int) bool { return res[i] < res[j] })
	// res = parallelSort(res, 20)
	out <- strings.Join(res, "_")
}

func ExecutePipeline(worker ...job) {
	var wg = &sync.WaitGroup{}
	// Начальный входной канал
	in := make(chan interface{})
	for _, op := range worker {
		wg.Add(1)
		out := make(chan interface{})
		go func(worker job, in chan interface{}, out chan interface{}) {
			defer wg.Done()
			defer close(out)
			worker(in, out)
		}(op, in, out)
		// Входной канал для следующей горутины — это выходной канал текущей
		in = out
	}
	// close(in)
	wg.Wait()
}

func main() {
	in := make(chan interface{})
	out := make(chan interface{})
	res := make(chan interface{})
	ans := make(chan interface{})
	start := time.Now()
	go SingleHash(in, out)
	go MultiHash(out, res)
	go CombineResults(res, ans)
	in <- "0"
	in <- "1"
	close(in)
	// in <- "Как  дела?"
	v := (<-ans).(string)
	fmt.Println(v)
	fmt.Println(time.Since(start))
}
