package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Station contains the station name as defined in the data input.
type Station string

// Temperature contains the float64 temperature as defined in the data input.
type Temperature float64

// Measurement is a struct containing a station name and it's recorded temperature.
// It is equivalent to a single line in the data input.
type Measurement struct {
	// Station contains the name of the station.
	Station Station
	// Temperature contains the temperature
	Temperature Temperature
}

// Info contains the aggregate data of a station's recorded temperatures.
type Info struct {
	// Count is the number of occurences of this station's recorded temperature.
	Count uint
	// Sum is the sum of all recorded temperatures.
	Sum Temperature
	// Min is the lowest recorded temperature for this station.
	Min Temperature
	// Max is the highest recorded temperature for this station.
	Max Temperature
}

// NewInfo returns an initial Info that represents a single observed temperature.
// This should be used when a new station's temperature is observed.
func NewInfo(temperature Temperature) *Info {
	return &Info{
		Count: 1,
		Sum:   temperature,
		Min:   temperature,
		Max:   temperature,
	}
}

// Update updates the aggregate data fields based on a new observed temperature.
func (i *Info) Update(temperature Temperature) {
	i.Count++
	i.Sum += temperature
	i.Min = min(i.Min, temperature)
	i.Max = max(i.Max, temperature)
}

// StationInfo is map mapping station names to their aggregate data.
type StationInfo map[Station]*Info

// AddInfo updates a station's Info based on a new Measurement.
func (s StationInfo) AddInfo(m Measurement) {
	if info, ok := s[m.Station]; !ok {
		s[m.Station] = NewInfo(m.Temperature)
	} else {
		info.Update(m.Temperature)
	}
}

// GenerateReport returns the entire output of StationInfo in the correct format as a string.
// The string returned should be the entire output required to write to stdout.
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

// ParseMeasurement parses a line from the input data and returns a Measurement.
func ParseMeasurement(line string) Measurement {
	s, tStr, _ := strings.Cut(line, ";")
	t, _ := strconv.ParseFloat(tStr, 64)
	return Measurement{
		Station:     Station(s),
		Temperature: Temperature(t),
	}
}

// StationReport formats a specific station's name, min, mean, and max, with 1 decimal precision.
// This adheres to a single line in the output to stdout.
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
