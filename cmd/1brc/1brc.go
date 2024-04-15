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

// Temperature contains an 16bit integer representation of temperature.
//
// |Temperature % 10| is the fractional component of the temperature.
//
// Temperature / 10 is the whole number component of the temperature.
//
// Temperature has a range of [-999, 999] that fits in an int16.
type Temperature int16

// TemperatureSum is accumulative int64 that holds the sum of Temperature.
//
// Since there are 1 billion expected temperatures, where Temperature is within [-999, 999],
// the sum will fit in an int64.
type TemperatureSum int64

// String returns a string representation of the Temperature,
// with as a number with 1 decimal placing.
func TemperatureString(t Temperature) string {
	isNegative := t < 0
	whole := t / 10
	frac := t % 10
	if isNegative {
		if whole == 0 {
			return fmt.Sprintf("-%d.%d", whole, -frac)
		} else if frac < 0 {
			return fmt.Sprintf("%d.%d", whole, -frac)
		} else {
			return fmt.Sprintf("%d.%d", whole, frac)
		}
	}
	return fmt.Sprintf("%d.%d", whole, frac)
}

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
	Sum TemperatureSum
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
		Sum:   TemperatureSum(temperature),
		Min:   temperature,
		Max:   temperature,
	}
}

// Update updates the aggregate data fields based on a new observed temperature.
func (i *Info) Update(temperature Temperature) {
	i.Count++
	i.Sum += TemperatureSum(temperature)
	i.Min = min(i.Min, temperature)
	i.Max = max(i.Max, temperature)
}

func (i Info) Mean() Temperature {
	roundingDigit := (int(i.Sum) * 10 / int(i.Count)) % 10
	if roundingDigit < 0 {
		roundingDigit = -roundingDigit
	}
	mean := int(i.Sum) / int(i.Count)
	if roundingDigit >= 5 {
		if mean < 0 {
			mean--
		} else {
			mean++
		}
	}
	return Temperature(mean)
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
	isNegative := tStr[0] == '-'
	whole, fraction, _ := strings.Cut(tStr, ".")
	w, _ := strconv.Atoi(whole)
	f, _ := strconv.Atoi(fraction)
	t := w * 10
	if isNegative {
		t -= f
	} else {
		t += f
	}
	return Measurement{
		Station:     Station(s),
		Temperature: Temperature(t),
	}
}

// StationReport formats a specific station's name, min, mean, and max, with 1 decimal precision.
// This adheres to a single line in the output to stdout.
func StationReport(station Station, info Info) string {
	return fmt.Sprintf(
		"%s;%s;%s;%s",
		station,
		TemperatureString(info.Min),
		TemperatureString(info.Mean()),
		TemperatureString(info.Max),
	)
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
