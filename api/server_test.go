package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFilterYearHandler(t *testing.T) {
	// Create a fake request
	req, err := http.NewRequest("GET", "/get_year?year=2023", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create a fake server
	srv := &Server{
		HolidayData: Holiday{
			EnglandAndWales: Division{
				Name: "england-and-wales",
				Events: []Event{
					{Title: "Event 1", Date: "2022-01-01", Bunting: true},
					{Title: "Event 2", Date: "2023-02-01", Bunting: false},
				},
			},
		},
	}

	srv.filterYearHandler(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %d, expected %d", status, http.StatusOK)
	}

	// Check the response body
	expected := `{"england-and-wales":{"division":"england-and-wales","events":[{"bunting":false,"date":"2023-02-01","notes":"","title":"Event 2"}]}}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v, expected %v", rr.Body.String(), expected)
	}
}

func TestFilterDivisionBuntingHandler(t *testing.T) {
	// Create a fake request
	req, err := http.NewRequest("GET", "/get_division_bunting?division=england-and-wales&bunting=false", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create a fake server
	srv := &Server{
		HolidayData: Holiday{
			EnglandAndWales: Division{
				Name: "england-and-wales",
				Events: []Event{
					{Title: "Event 1", Date: "2022-01-01", Bunting: true},
					{Title: "Event 2", Date: "2023-02-01", Bunting: false},
				},
			},
		},
	}

	srv.filterDivisionBuntingHandler(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %d, expected %d", status, http.StatusOK)
	}

	// Check the response body
	expected := `{"england-and-wales":{"division":"england-and-wales","events":[{"bunting":false,"date":"2023-02-01","notes":"","title":"Event 2"}]}}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v, expected %v", rr.Body.String(), expected)
	}
}

func TestFilterFilterYearWithRestrictedEventHandler(t *testing.T) {
	// Create a fake request
	req, err := http.NewRequest("GET", "get_year_less_event?year=2023", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create a fake server
	srv := &Server{
		HolidayData: Holiday{
			EnglandAndWales: Division{
				Name: "england-and-wales",
				Events: []Event{
					{Title: "Event 1", Date: "2022-01-01", Bunting: true, Notes: "Note 1"},
					{Title: "Event 2", Date: "2023-02-01", Bunting: false, Notes: "Note 2"},
				},
			},
		},
	}

	srv.filterYearWithRestrictedEventHandler(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %d, expected %d", status, http.StatusOK)
	}

	// Check the response body
	expected := `{"england-and-wales":{"division":"england-and-wales","events":[{"date":"2023-02-01","title":"Event 2"}]}}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v, expected %v", rr.Body.String(), expected)
	}
}
