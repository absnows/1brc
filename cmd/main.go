package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Temperture struct {
	count int64
	temps []float64
	sum   float64
}

func main() {
	f := flag.String("f", "foo", "Filepath to read")
	flag.Parse()
	fmt.Println(evaluate(*f))
}

func evaluate(path string) string {
	information, err := speedRead(path)
	if err != nil {
		panic(err)
	}

	ch := make(chan string, len(information))
	var result []string
	var wg sync.WaitGroup

	for c, temperture := range information {
		wg.Add(1)
		go func(c string, temperture Temperture) {
			sort.Float64s(temperture.temps)
			avg := temperture.sum / float64(temperture.count)
			ch <- fmt.Sprintf(
				"%s=%.1f/%.1f/%.1f",
				c,
				temperture.temps[0],
				avg,
				temperture.temps[len(temperture.temps)-1])
			wg.Done()
		}(c, temperture)
	}

	wg.Wait()
	for i := 0; i < len(information); i++ {
		result = append(result, <-ch)
	}

	close(ch)
	sort.Strings(result)
	return strings.Join(result, ", ")

}

func speedRead(path string) (map[string]Temperture, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	tempertures := make(chan map[string]Temperture, 10)
	chunkStream := make(chan []byte, 15)
	chunkSize := 64 * 1024 * 1024

	var wg sync.WaitGroup

	for i := 0; i < runtime.NumCPU()-1; i++ {
		wg.Add(1)
		go func() {
			for chunk := range chunkStream {
				readChunk(chunk, tempertures)
			}
			wg.Done()
		}()
	}

	go func() {
		buf := make([]byte, chunkSize)
		leftover := make([]byte, 0, chunkSize)

		for {
			total, err := f.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				panic(err)
			}

			buf = buf[:total]

			send := make([]byte, total)
			copy(send, buf)

			lastNewLineIndex := bytes.LastIndex(buf, []byte{'\n'})

			send = append(leftover, buf[:lastNewLineIndex+1]...)
			leftover = make([]byte, len(buf[lastNewLineIndex+1:]))
			copy(leftover, buf[lastNewLineIndex+1:])

			chunkStream <- send
		}

		close(chunkStream)

		wg.Wait()
		close(tempertures)
	}()

	response := make(map[string]Temperture)
	for t := range tempertures {
		for city, temperture := range t {
			response[city] = temperture
		}
	}

	return response, nil
}

func readChunk(buf []byte, result chan<- map[string]Temperture) {
	send := make(map[string]Temperture)

	var start int
	var city string

	stringBuf := string(buf)
	for index, char := range stringBuf {
		switch char {
		case ';':
			city = stringBuf[start:index]
			start = index + 1
		case '\n':
			if (index-start) > 1 && len(city) != 0 {
				temp := convert(stringBuf[start:index])
				start = index + 1

				if val, ok := send[city]; ok {
					val.count++
					val.sum += temp
					val.temps = append(val.temps, temp)
					send[city] = val
				} else {
					send[city] = Temperture{
						count: 1,
						temps: []float64{temp},
						sum:   temp,
					}
				}

				city = ""
			}
		}
	}

	result <- send
}

func convert(s string) float64 {
	output, _ := strconv.ParseFloat(s, 64)
	return output
}
