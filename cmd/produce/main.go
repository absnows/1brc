package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand/v2"
	"os"
	"strings"
)

func main() {
	size := flag.Int("s", 10, "Size of file")
	output := flag.String("o", "mensurement", "File output name")
	flag.Parse()

	file, err := os.OpenFile("./data/weather_stations.csv", os.O_RDWR, 0600)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	rb, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	stations := strings.Split(string(rb), "\n")

	generateFile(*size, stations, *output)

}

func generateFile(size int, stations []string, out string) {
	f, err := os.Create("./data/" + out + ".txt")
	if err != nil {
		log.Fatal(err)
		return
	}

	defer f.Close()

	for i := 0; i < size; i++ {
		s := strings.Split(getSomeStation(stations), ";")[0]
		tempeture := getRandomFloat(-99.9, 99.9)
		line := fmt.Sprintf("%s;%f", s, tempeture)
		_, err := f.WriteString(line + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getSomeStation(s []string) string {
	sts := s[2:]
	return sts[rand.IntN(44690)]
}

func getRandomFloat(min, max float64) float64 {
	n := min + rand.Float64()*(max-min)
	return math.Round(n*100) / 100
}
