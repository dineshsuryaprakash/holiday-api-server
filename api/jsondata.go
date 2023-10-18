package api

import (
	"encoding/json"
	"net/http"

	"k8s.io/klog/v2"
)

type Event struct {
	Title   string `json:"title"`
	Date    string `json:"date"`
	Notes   string `json:"notes"`
	Bunting bool   `json:"bunting"`
}

type Division struct {
	Name   string  `json:"division,omitempty"`
	Events []Event `json:"events,omitempty"`
}

type Holiday struct {
	EnglandAndWales Division `json:"england-and-wales"`
	Scotland        Division `json:"scotland"`
	NorthernIreland Division `json:"northern-ireland"`
}

func ParseHolidayData(url string, holiday *Holiday) error {
	resp, err := http.Get(url)
	if err != nil {
		klog.Errorf("Error fetching data")
		return err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(holiday); err != nil {
		klog.Errorf("Error decoding JSON")
		return err
	}
	return nil
}
