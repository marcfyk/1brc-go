package main

import (
	"reflect"
	"testing"
)

func TestTemperatureString(t *testing.T) {
	tests := []struct {
		name     string
		t        Temperature
		expected string
	}{
		{
			name:     "zero temperature",
			t:        0,
			expected: "0.0",
		},
		{
			name:     "temperature with no decimal place",
			t:        200,
			expected: "20.0",
		},
		{
			name:     "temperature with decimal place",
			t:        207,
			expected: "20.7",
		},
		{
			name:     "negative temperature",
			t:        -23,
			expected: "-2.3",
		},
		{
			name:     "zero whole number with non-zero fractional",
			t:        5,
			expected: "0.5",
		},
		{
			name:     "negative zero whole number with non-zero fractional",
			t:        -5,
			expected: "-0.5",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if s := TemperatureString(test.t); s != test.expected {
				t.Errorf("actual string: %s, expected string: %s", s, test.expected)
			}
		})
	}
}

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
				Min:   0,
				Max:   1000,
			},
			temperature: 100,
			expected: Info{
				Count: 3,
				Sum:   1100,
				Min:   0,
				Max:   1000,
			},
		},
		{
			name: "updating the min",
			info: Info{
				Count: 2,
				Sum:   1000,
				Min:   0,
				Max:   1000,
			},
			temperature: -1000,
			expected: Info{
				Count: 3,
				Sum:   0,
				Min:   -1000,
				Max:   1000,
			},
		},
		{
			name: "updating the max",
			info: Info{
				Count: 2,
				Sum:   1000,
				Min:   0,
				Max:   1000,
			},
			temperature: 2000,
			expected: Info{
				Count: 3,
				Sum:   3000,
				Min:   0,
				Max:   2000,
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
				Temperature: 10000,
			},
			expected: StationInfo{
				"Beijing": &Info{
					Count: 1,
					Sum:   10000,
					Max:   10000,
					Min:   10000,
				},
			},
		},
		{
			name: "subsequent adding of info",
			stationInfo: StationInfo{
				"Beijing": &Info{
					Count: 1,
					Sum:   1000,
					Max:   1000,
					Min:   1000,
				},
			},
			measurement: Measurement{
				Station:     "Beijing",
				Temperature: 505,
			},
			expected: StationInfo{
				"Beijing": &Info{
					Count: 2,
					Sum:   1505,
					Max:   1000,
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
			name: "formats values to 1 decimal point",
			stationInfo: StationInfo{
				"Beijing": &Info{
					Count: 10,
					Sum:   10000,
					Min:   0,
					Max:   5005,
				},
				"Ḩamīdīyeh": &Info{
					Count: 1,
					Sum:   0,
					Min:   0,
					Max:   0,
				},
				"New York": &Info{
					Count: 10,
					Sum:   200,
					Min:   0,
					Max:   400,
				},
			},
			report: "Beijing;0.0;100.0;500.5\nNew York;0.0;2.0;40.0\nḨamīdīyeh;0.0;0.0;0.0",
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

func TestParseMeasurement(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		measurement Measurement
	}{
		{
			name: "single word station",
			line: "Beijing;70.5",
			measurement: Measurement{
				Station:     "Beijing",
				Temperature: 705,
			},
		},
		{
			name: "station with spaces",
			line: "New York;70.5",
			measurement: Measurement{
				Station:     "New York",
				Temperature: 705,
			},
		},
		{
			name: "station with unicode",
			line: "Ḩamīdīyeh;70.5",
			measurement: Measurement{
				Station:     "Ḩamīdīyeh",
				Temperature: 705,
			},
		},
		{
			name: "negative temperature",
			line: "Paris;-10.2",
			measurement: Measurement{
				Station:     "Paris",
				Temperature: -102,
			},
		},
		{
			name: "zero temperature",
			line: "Berlin;0.0",
			measurement: Measurement{
				Station:     "Berlin",
				Temperature: 0,
			},
		},
		{
			name: "positive temperature with zero whole number component and non-zero fractional component",
			line: "Rome;0.7",
			measurement: Measurement{
				Station:     "Rome",
				Temperature: 7,
			},
		},
		{
			name: "negative temperature with zero whole number component and non-zero fractional component",
			line: "Rome;-0.7",
			measurement: Measurement{
				Station:     "Rome",
				Temperature: -7,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := ParseMeasurement(test.line)
			if m != test.measurement {
				t.Errorf("actual measurement: %+v, expected measurement: %+v", m, test.measurement)
			}
		})
	}
}

func TestStationReport(t *testing.T) {
	tests := []struct {
		name    string
		station Station
		info    Info
		report  string
	}{
		{
			name:    "single word station with multiple counts",
			station: "Beijing",
			info: Info{
				Count: 5,
				Sum:   1005,
				Min:   101,
				Max:   300,
			},
			report: "Beijing;10.1;20.1;30.0",
		},
		{
			name:    "station with spaces",
			station: "New York",
			info: Info{
				Count: 5,
				Sum:   1005,
				Min:   101,
				Max:   300,
			},
			report: "New York;10.1;20.1;30.0",
		},
		{
			name:    "station with unicode",
			station: "Ḩamīdīyeh",
			info: Info{
				Count: 5,
				Sum:   1005,
				Min:   101,
				Max:   300,
			},
			report: "Ḩamīdīyeh;10.1;20.1;30.0",
		},
		{
			name:    "single count",
			station: "Beijing",
			info: Info{
				Count: 1,
				Sum:   100,
				Min:   100,
				Max:   100,
			},
			report: "Beijing;10.0;10.0;10.0",
		},
		{
			name:    "negative values",
			station: "Beijing",
			info: Info{
				Count: 4,
				Sum:   -8000,
				Min:   -3005,
				Max:   -1007,
			},
			report: "Beijing;-300.5;-200.0;-100.7",
		},
		{
			name:    "rounded down mean",
			station: "Madrid",
			info: Info{
				Count: 7,
				Sum:   107,
				Min:   -7,
				Max:   120,
			},
			report: "Madrid;-0.7;1.5;12.0",
		},
		{
			name:    "negative rounded down mean",
			station: "Madrid",
			info: Info{
				Count: 7,
				Sum:   -107,
				Min:   -120,
				Max:   7,
			},
			report: "Madrid;-12.0;-1.5;0.7",
		},
		{
			name:    "rounded up mean",
			station: "Barcelona",
			info: Info{
				Count: 7,
				Sum:   102,
				Min:   -10,
				Max:   200,
			},
			report: "Barcelona;-1.0;1.5;20.0",
		},
		{
			name:    "negative rounded up mean",
			station: "Barcelona",
			info: Info{
				Count: 7,
				Sum:   -102,
				Min:   -200,
				Max:   10,
			},
			report: "Barcelona;-20.0;-1.5;1.0",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if report := StationReport(test.station, test.info); report != test.report {
				t.Errorf("actual report: %+v, expected report: %+v", report, test.report)
			}
		})
	}
}
