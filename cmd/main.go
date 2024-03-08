package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func main() {
	f := flag.String("f", "foo", "Filepath to read")
	flag.Parse()
	fmt.Println(evaluate(*f))
}

func evaluate(path string) string {
	information, err := readFileEachLine(path)
	if err != nil {
		panic(err)
	}

	ch := make(chan string, len(information))
	var result []string
	var wg sync.WaitGroup

	for c, temps := range information {
		wg.Add(1)
		go func(c string, temps []float64) {
			sort.Float64s(temps)

			var avg float64
			for _, temp := range temps {
				avg += temp
			}

			avg = avg / float64(len(temps))
			avg = math.Ceil(avg*10) / 10

			ch <- fmt.Sprintf("%s=%.1f/%.1f/%.1f", c, temps[0], avg, temps[len(temps)-1])
			wg.Done()
		}(c, temps)
	}

	wg.Wait()
	for i := 0; i < len(information); i++ {
		result = append(result, <-ch)
	}

	close(ch)
	sort.Strings(result)
	return strings.Join(result, ", ")

}

func readFileEachLine(path string) (map[string][]float64, error) {
	m := make(map[string][]float64)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		text := scanner.Text()

		index := strings.Index(text, ";")
		c := text[:index]
		t := convert(text[index+1:])

		if _, ok := m[c]; ok {
			m[c] = append(m[c], t)
		} else {
			m[c] = []float64{t}
		}
	}

	return m, nil
}

func convert(s string) float64 {
	output, _ := strconv.ParseFloat(s, 64)
	return output
}
