package main

import (
	"reflect"
	"testing"
)

func TestNewInfo(t *testing.T) {
	i := NewInfo(15)
	expected := Info{
		Count: 1,
		Sum:   15,
		Min:   15,
		Max:   15,
	}
	if *i != expected {
		t.Errorf("actual info: %+v, expected info: %+v", *i, expected)
	}
}

func TestInfoUpdate(t *testing.T) {
	tests := []struct {
		name        string
		info        Info
		temperature Temperature
		expected    Info
	}{
		{
			name: "updating the sum",
			info: Info{
				Count: 2,
				Sum:   1000,
				Min:   1,
				Max:   999,
			},
			temperature: 100,
			expected: Info{
				Count: 3,
				Sum:   1100,
				Min:   1,
				Max:   999,
			},
		},
		{
			name: "updating the min",
			info: Info{
				Count: 2,
				Sum:   1000,
				Min:   1,
				Max:   999,
			},
			temperature: -999,
			expected: Info{
				Count: 3,
				Sum:   1,
				Min:   -999,
				Max:   999,
			},
		},
		{
			name: "updating the max",
			info: Info{
				Count: 2,
				Sum:   999,
				Min:   0,
				Max:   999,
			},
			temperature: 550,
			expected: Info{
				Count: 3,
				Sum:   1549,
				Min:   0,
				Max:   999,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.info.Update(test.temperature)
			if !reflect.DeepEqual(test.info, test.expected) {
				t.Error()
			}
		})
	}
}

func TestStationInfoAddInfo(t *testing.T) {
	tests := []struct {
		name        string
		stationInfo StationInfo
		measurement Measurement
		expected    StationInfo
	}{
		{
			name:        "first time adding info",
			stationInfo: StationInfo{},
			measurement: Measurement{
				Station:     "Beijing",
				Temperature: 999,
			},
			expected: StationInfo{
				"Beijing": &Info{
					Count: 1,
					Sum:   999,
					Min:   999,
					Max:   999,
				},
			},
		},
		{
			name: "subsequent adding of info",
			stationInfo: StationInfo{
				"Beijing": &Info{
					Count: 1,
					Sum:   999,
					Min:   999,
					Max:   999,
				},
			},
			measurement: Measurement{
				Station:     "Beijing",
				Temperature: 505,
			},
			expected: StationInfo{
				"Beijing": &Info{
					Count: 2,
					Sum:   1504,
					Max:   999,
					Min:   505,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.stationInfo.AddInfo(test.measurement)
			if !reflect.DeepEqual(test.stationInfo, test.expected) {
				t.Errorf("actual station info: %+v, expected station info: %+v", test.stationInfo, test.expected)
			}
		})
	}
}

func TestStationInfoGenerateReport(t *testing.T) {
	tests := []struct {
		name        string
		stationInfo StationInfo
		report      string
	}{
		{
			name: "zero temperature",
			stationInfo: StationInfo{
				"London": {
					Count: 1,
					Sum:   0,
					Min:   0,
					Max:   0,
				},
			},
			report: "London;0.0;0.0;0.0",
		},
		{
			name: "temperature with no decimal place",
			stationInfo: StationInfo{
				"Liverpool": {
					Count: 1,
					Sum:   200,
					Min:   200,
					Max:   200,
				},
			},
			report: "Liverpool;20.0;20.0;20.0",
		},
		{
			name: "temperature with decimal place",
			stationInfo: StationInfo{
				"Manchester": {
					Count: 1,
					Sum:   207,
					Min:   207,
					Max:   207,
				},
			},
			report: "Manchester;20.7;20.7;20.7",
		},
		{
			name: "negative temperature",
			stationInfo: StationInfo{
				"Whales": {
					Count: 1,
					Sum:   -23,
					Min:   -23,
					Max:   -23,
				},
			},
			report: "Whales;-2.3;-2.3;-2.3",
		},
		{
			name: "zero whole number with non-zero fracional",
			stationInfo: StationInfo{
				"Venice": {
					Count: 1,
					Sum:   5,
					Min:   5,
					Max:   5,
				},
			},
			report: "Venice;0.5;0.5;0.5",
		},
		{
			name: "negative zero whole number with non-zero fracional",
			stationInfo: StationInfo{
				"Venice": {
					Count: 1,
					Sum:   -5,
					Min:   -5,
					Max:   -5,
				},
			},
			report: "Venice;-0.5;-0.5;-0.5",
		},
		{
			name: "single word station with multiple counts",
			stationInfo: StationInfo{
				"Florence": {
					Count: 5,
					Sum:   1005,
					Min:   101,
					Max:   300,
				},
			},
			report: "Florence;10.1;20.1;30.0",
		},
		{
			name: "station with spaces",
			stationInfo: StationInfo{
				"New York": {
					Count: 5,
					Sum:   1005,
					Min:   101,
					Max:   300,
				},
			},
			report: "New York;10.1;20.1;30.0",
		},
		{
			name: "station with unicode",
			stationInfo: StationInfo{
				"Ḩamīdīyeh": {
					Count: 5,
					Sum:   1005,
					Min:   101,
					Max:   300,
				},
			},
			report: "Ḩamīdīyeh;10.1;20.1;30.0",
		},
		{
			name: "rounded down mean",
			stationInfo: StationInfo{
				"Barcelona": {
					Count: 7,
					Sum:   107,
					Min:   -7,
					Max:   120,
				},
			},
			report: "Barcelona;-0.7;1.5;12.0",
		},
		{
			name: "negative rounded down mean",
			stationInfo: StationInfo{
				"Barcelona": {
					Count: 7,
					Sum:   -107,
					Min:   -120,
					Max:   7,
				},
			},
			report: "Barcelona;-12.0;-1.5;0.7",
		},
		{
			name: "rounded up mean",
			stationInfo: StationInfo{
				"Madrid": {
					Count: 7,
					Sum:   102,
					Min:   -10,
					Max:   200,
				},
			},
			report: "Madrid;-1.0;1.5;20.0",
		},
		{
			name: "negative rounded up mean",
			stationInfo: StationInfo{
				"Madrid": {
					Count: 7,
					Sum:   -102,
					Min:   -200,
					Max:   10,
				},
			},
			report: "Madrid;-20.0;-1.5;1.0",
		},
		{
			name: "formats multiple values separated by newlines sorted alphabetically",
			stationInfo: StationInfo{
				"Beijing": {
					Count: 10,
					Sum:   9990,
					Min:   0,
					Max:   999,
				},
				"Ḩamīdīyeh": {
					Count: 1,
					Sum:   0,
					Min:   0,
					Max:   0,
				},
				"New York": {
					Count: 10,
					Sum:   200,
					Min:   0,
					Max:   400,
				},
			},
			report: "Beijing;0.0;99.9;99.9\nNew York;0.0;2.0;40.0\nḨamīdīyeh;0.0;0.0;0.0",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if report := test.stationInfo.GenerateReport(); report != test.report {
				t.Errorf("actual report: %s, expected report: %s", report, test.report)
			}
		})
	}
}
