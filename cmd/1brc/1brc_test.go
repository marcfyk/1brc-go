package main

import (
	"reflect"
	"testing"
)

func TestNewInfo(t *testing.T) {
	i := NewInfo(1.5)
	expected := Info{
		Count: 1,
		Sum:   1.5,
		Min:   1.5,
		Max:   1.5,
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
				Sum:   100,
				Min:   0,
				Max:   100,
			},
			temperature: 10,
			expected: Info{
				Count: 3,
				Sum:   110,
				Min:   0,
				Max:   100,
			},
		},
		{
			name: "updating the min",
			info: Info{
				Count: 2,
				Sum:   100,
				Min:   0,
				Max:   100,
			},
			temperature: -100,
			expected: Info{
				Count: 3,
				Sum:   0,
				Min:   -100,
				Max:   100,
			},
		},
		{
			name: "updating the max",
			info: Info{
				Count: 2,
				Sum:   100,
				Min:   0,
				Max:   100,
			},
			temperature: 200,
			expected: Info{
				Count: 3,
				Sum:   300,
				Min:   0,
				Max:   200,
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
				Temperature: 100,
			},
			expected: StationInfo{
				"Beijing": &Info{
					Count: 1,
					Sum:   100,
					Max:   100,
					Min:   100,
				},
			},
		},
		{
			name: "subsequent adding of info",
			stationInfo: StationInfo{
				"Beijing": &Info{
					Count: 1,
					Sum:   100,
					Max:   100,
					Min:   100,
				},
			},
			measurement: Measurement{
				Station:     "Beijing",
				Temperature: 50.5,
			},
			expected: StationInfo{
				"Beijing": &Info{
					Count: 2,
					Sum:   150.5,
					Max:   100,
					Min:   50.5,
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
					Sum:   1000,
					Min:   0,
					Max:   500.5,
				},
				"Ḩamīdīyeh": &Info{
					Count: 1,
					Sum:   0,
					Min:   0,
					Max:   0,
				},
				"New York": &Info{
					Count: 10,
					Sum:   20,
					Min:   0,
					Max:   40,
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
				Temperature: 70.5,
			},
		},
		{
			name: "station with spaces",
			line: "New York;70.5",
			measurement: Measurement{
				Station:     "New York",
				Temperature: 70.5,
			},
		},
		{
			name: "station with unicode",
			line: "Ḩamīdīyeh;70.5",
			measurement: Measurement{
				Station:     "Ḩamīdīyeh",
				Temperature: 70.5,
			},
		},
		{
			name: "negative temperature",
			line: "Paris;-10.2",
			measurement: Measurement{
				Station:     "Paris",
				Temperature: -10.2,
			},
		},
		{
			name: "zero temperature",
			line: "Berlin;0.0",
			measurement: Measurement{
				Station:     "Berlin",
				Temperature: 0.0,
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
				Sum:   100.5,
				Min:   10.1,
				Max:   30,
			},
			report: "Beijing;10.1;20.1;30.0",
		},
		{
			name:    "station with spaces",
			station: "New York",
			info: Info{
				Count: 5,
				Sum:   100.5,
				Min:   10.1,
				Max:   30,
			},
			report: "New York;10.1;20.1;30.0",
		},
		{
			name:    "station with unicode",
			station: "Ḩamīdīyeh",
			info: Info{
				Count: 5,
				Sum:   100.5,
				Min:   10.1,
				Max:   30,
			},
			report: "Ḩamīdīyeh;10.1;20.1;30.0",
		},
		{
			name:    "single count",
			station: "Beijing",
			info: Info{
				Count: 1,
				Sum:   10,
				Min:   10,
				Max:   10,
			},
			report: "Beijing;10.0;10.0;10.0",
		},
		{
			name:    "negative values",
			station: "Beijing",
			info: Info{
				Count: 4,
				Sum:   -800,
				Min:   -300,
				Max:   -100,
			},
			report: "Beijing;-300.0;-200.0;-100.0",
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
