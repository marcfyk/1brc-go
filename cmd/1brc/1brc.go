package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"sort"
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
	Count uint32
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
	sum, count := int32(i.Sum), int32(i.Count)
	roundingDigit := (sum * 10 / count) % 10
	if roundingDigit < 0 {
		roundingDigit = -roundingDigit
	}
	mean := sum / count
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

// AppendTemperatureToBuffer writes temperature to a byte buffer at the given cursor
// and returns the new cursor position.
func AppendTemperatureToBuffer(buffer []byte, cursor int, t Temperature) int {
	n := int16(t)
	if isNegative := n < 0; isNegative {
		buffer[cursor] = '-'
		cursor++
		n = -n
	}
	if d := n / 100; d > 0 {
		buffer[cursor] = byte(d + '0')
		cursor++
	}
	buffer[cursor] = byte(n/10%10 + '0')
	cursor++
	buffer[cursor] = '.'
	cursor++
	buffer[cursor] = byte(n%10 + '0')
	cursor++
	return cursor
}

// GenerateReport returns the entire output of StationInfo in the ordered by station name as a string.
// The string returned should be the entire output required to write to stdout.
func (s StationInfo) GenerateReport() string {
	stations := make([]Station, 0, len(s))
	bufferSize := len(s) * (1 + 5 + 1 + 5 + 1 + 5 + 1) // each line has <station;digit;digit;digit>
	for station := range s {
		bufferSize += len([]byte(station))
		stations = append(stations, station)
	}
	sort.Slice(stations, func(i, j int) bool { return stations[i] < stations[j] })
	b := make([]byte, bufferSize)
	cursor := 0
	for i, station := range stations {
		for _, c := range []byte(station) {
			b[cursor] = c
			cursor++
		}
		b[cursor] = ';'
		cursor++
		cursor = AppendTemperatureToBuffer(b, cursor, s[station].Min)
		b[cursor] = ';'
		cursor++
		cursor = AppendTemperatureToBuffer(b, cursor, s[station].Mean())
		b[cursor] = ';'
		cursor++
		cursor = AppendTemperatureToBuffer(b, cursor, s[station].Max)
		if i < len(stations)-1 {
			b[cursor] = '\n'
			cursor++
		}
	}
	return string(b[:cursor])
}

func main() {
	fp := os.Args[1]
	f, _ := os.Open(fp)
	defer f.Close()
	b := make([]byte, 1<<25)
	cursor := 0
	stationInfo := make(StationInfo)
	for {
		n, err := f.Read(b[cursor:])
		if err != nil && !errors.Is(err, io.EOF) {
			break
		}
		if n == 0 {
			break
		}
		chunk := b[:cursor+n]
		index := bytes.LastIndexByte(chunk, '\n')
		if index == -1 {
			break
		}
		unprocessed := chunk[index+1:]
		chunk = chunk[:index]
		for {
			semicolon := slices.Index(chunk, ';')
			if semicolon == -1 {
				break
			}
			newline := slices.Index(chunk, '\n')
			if newline == -1 {
				newline = len(chunk)
			}
			var t Temperature
			switch semicolon {
			case newline - 4:
				t = Temperature(chunk[newline-3]-'0')*10 + Temperature(chunk[newline-1]-'0')
			case newline - 5:
				if chunk[newline-4] == '-' {
					t = -(Temperature(chunk[newline-3]-'0')*10 + Temperature(chunk[newline-1]-'0'))
				} else {
					t = Temperature(chunk[newline-4]-'0')*100 + Temperature(chunk[newline-3]-'0')*10 + Temperature(chunk[newline-1]-'0')
				}
			case newline - 6:
				t = -(Temperature(chunk[newline-4]-'0')*100 + Temperature(chunk[newline-3]-'0')*10 + Temperature(chunk[newline-1]-'0'))
			}
			m := Measurement{
				Station:     Station(chunk[:semicolon]),
				Temperature: t,
			}
			stationInfo.AddInfo(m)
			if newline == len(chunk) {
				break
			}
			chunk = chunk[newline+1:]
		}
		cursor = copy(b, unprocessed)
	}
	fmt.Println(stationInfo.GenerateReport())
}
