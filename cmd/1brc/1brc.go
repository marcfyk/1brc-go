package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Station string
type Temperature float64

type Measurement struct {
	Station     Station
	Temperature Temperature
}

type Info struct {
	Count uint
	Sum   Temperature
	Min   Temperature
	Max   Temperature
}

func NewInfo(temperature Temperature) *Info {
	return &Info{
		Count: 1,
		Sum:   temperature,
		Min:   temperature,
		Max:   temperature,
	}
}

func (i *Info) Update(temperature Temperature) {
	i.Count++
	i.Sum += temperature
	i.Min = min(i.Min, temperature)
	i.Max = max(i.Max, temperature)
}

type StationInfo map[Station]*Info

func (s StationInfo) AddInfo(m Measurement) {
	if info, ok := s[m.Station]; !ok {
		s[m.Station] = NewInfo(m.Temperature)
	} else {
		info.Update(m.Temperature)
	}
}

func (s StationInfo) GenerateReport() string {
	stations := make([]Station, 0, len(s))
	for station := range s {
		stations = append(stations, station)
	}
	sort.Slice(stations, func(i, j int) bool { return stations[i] < stations[j] })
	output := make([]string, len(stations))
	for i, station := range stations {
		output[i] = StationReport(station, *s[station])
	}
	return strings.Join(output, "\n")
}

func ParseMeasurement(line string) Measurement {
	tokens := strings.Split(line, ";")
	station := tokens[0]
	temperature, _ := strconv.ParseFloat(tokens[1], 64)
	return Measurement{
		Station:     Station(station),
		Temperature: Temperature(temperature),
	}
}

func StationReport(station Station, info Info) string {
	mean := info.Sum / Temperature(info.Count)
	return fmt.Sprintf("%s;%.1f;%.1f;%.1f", station, info.Min, mean, info.Max)
}

func main() {
	fp := os.Args[1]
	f, _ := os.Open(fp)
	defer f.Close()
	stationInfo := make(StationInfo)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		m := ParseMeasurement(line)
		stationInfo.AddInfo(m)
	}
	fmt.Println(stationInfo.GenerateReport())
}
