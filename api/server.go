package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"k8s.io/klog/v2"
)

type Server struct {
	WebServer       *http.Server
	JsonEndpointURL string
	HolidayData     Holiday
}

// NewServer creates a new instance of the Server struct
func NewServer(listenAddr string, BankHolidayJsonEndpointURL string) (*Server, error) {
	holiday := Holiday{}
	if err := ParseHolidayData(BankHolidayJsonEndpointURL, &holiday); err != nil {
		klog.Errorf("Parsing holiday data failed")
		return nil, err
	}

	klog.Infof("Server created with JSON endpoint URL: %s", BankHolidayJsonEndpointURL)

	return &Server{
		WebServer: &http.Server{
			Addr: listenAddr,
		},
		JsonEndpointURL: BankHolidayJsonEndpointURL,
		HolidayData:     holiday,
	}, nil
}

// Initialize the server by registering all the handle functions
func (s *Server) Initialize() {
	// USEAGE: GET http://localhost:8080/get_year?year=2023
	http.HandleFunc("/get_year", s.filterYearHandler)
	// USEAGE: GET http://localhost:8080/get_division_bunting?division=england-and-wales&bunting=true
	http.HandleFunc("/get_division_bunting", s.filterDivisionBuntingHandler)
	// USEAGE: GET http://localhost:8080/get_year_less_event?year=2023
	http.HandleFunc("/get_year_less_event", s.filterYearWithRestrictedEventHandler)
}

// Run listens on the TCP network address and serves incoming HTTP requests.
func (s *Server) Run() {
	go func() {
		klog.Infof("Server listening on : %s", s.WebServer.Addr)
		if err := s.WebServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			klog.Errorf("Error Listening and serving Connetctions: %s", err)
		}
	}()
}

// Shutdown performs gracefully shut down the server
func (s *Server) Shutdown() error {
	klog.Infof("Initializing server shudown")

	return s.WebServer.Shutdown(context.Background())
}

func findstring(target string, slice []string) bool {
	for _, element := range slice {
		if element == target {
			return true
		}
	}
	return false
}

func removeFields(data *interface{}, fieldsToRemove []string) {
	switch (*data).(type) {
	case map[string]interface{}:
		valuesMap := (*data).(map[string]interface{})
		// Iterate through the map
		for key, val := range valuesMap {
			if findstring(key, fieldsToRemove) {
				delete(valuesMap, key)
			} else {
				removeFields(&val, fieldsToRemove)
			}
		}
	case []interface{}:
		for _, val := range (*data).([]interface{}) {
			removeFields(&val, fieldsToRemove)

		}
	}
}

func cleanJson(jsonStr []byte) ([]byte, error) {
	// Unmarshal the JSON into a map
	var data map[string]interface{}
	if err := json.Unmarshal(jsonStr, &data); err != nil {
		return nil, err
	}

	// Check if "england-and-wales" should be removed
	englandAndWalesDivision := data["england-and-wales"].(map[string]interface{})
	if name, ok := englandAndWalesDivision["division"].(string); !ok || name == "" {
		delete(data, "england-and-wales")
	}

	// Check if "scotland" should be removed
	scotlandDivision := data["scotland"].(map[string]interface{})
	if name, ok := scotlandDivision["division"].(string); !ok || name == "" {
		delete(data, "scotland")
	}

	// Check if "northern-ireland" should be removed
	northernirelandDivision := data["northern-ireland"].(map[string]interface{})
	if name, ok := northernirelandDivision["division"].(string); !ok || name == "" {
		delete(data, "northern-ireland")
	}

	// Marshal the modified map back to JSON
	modifiedJSON, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return modifiedJSON, nil
}

func writeHttpResponse(w http.ResponseWriter, toJson interface{}) {
	// Marshal the modified data back into a JSON string
	JSON, err := json.Marshal(toJson)
	if err != nil {
		error := fmt.Sprintf("WriteHttpResponse Error: %v", err)
		http.Error(w, error, http.StatusInternalServerError)
		return
	}

	JSON, err = cleanJson(JSON)
	if err != nil {
		error := fmt.Sprintf("WriteHttpResponse Error: %v", err)
		http.Error(w, error, http.StatusInternalServerError)
		return
	}

	// Write the JSON to the response writer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(JSON)
}

func filterEventsOnBunting(data []Event, expectedBuntingValue bool, divisionName string) Division {
	var filteredEvents []Event
	for _, event := range data {
		if event.Bunting == expectedBuntingValue {
			filteredEvents = append(filteredEvents, event)
		}
	}
	// Set the filtered events to the filteredDivision
	filteredDivision := Division{
		Name:   divisionName,
		Events: filteredEvents,
	}
	return filteredDivision
}

func filterYear(w http.ResponseWriter, r *http.Request, year int, HolidayData []Division) Holiday {

	filteredHoliday := Holiday{}

	// Loop through events in HolidayData
	for _, division := range HolidayData {
		for _, event := range division.Events {
			// Check if the event's year matches the specified year
			if strings.HasPrefix(event.Date, strconv.Itoa(year)) {
				// Append matching events to the respective division in the filteredHoliday
				switch division.Name {
				case "england-and-wales":
					if filteredHoliday.EnglandAndWales.Name == "" {
						filteredHoliday.EnglandAndWales.Name = "england-and-wales"
					}
					filteredHoliday.EnglandAndWales.Events = append(filteredHoliday.EnglandAndWales.Events, event)
				case "scotland":
					if filteredHoliday.Scotland.Name == "" {
						filteredHoliday.Scotland.Name = "scotland"
					}
					filteredHoliday.Scotland.Events = append(filteredHoliday.Scotland.Events, event)
				case "northern-ireland":
					if filteredHoliday.NorthernIreland.Name == "" {
						filteredHoliday.NorthernIreland.Name = "northern-ireland"
					}
					filteredHoliday.NorthernIreland.Events = append(filteredHoliday.NorthernIreland.Events, event)
				}
			}
		}
	}
	return filteredHoliday
}

// Handler function for the "/division_bunting" path
func (s *Server) filterDivisionBuntingHandler(w http.ResponseWriter, r *http.Request) {
	// Get the query parameters from the request
	divisionName := r.URL.Query().Get("division")
	buntingParam := r.URL.Query().Get("bunting")

	bunting, err := strconv.ParseBool(buntingParam)
	if err != nil {
		http.Error(w, "Invalid bunting parameter", http.StatusBadRequest)
		return
	}

	filteredHoliday := Holiday{}

	switch divisionName {
	case "england-and-wales":
		filteredHoliday.EnglandAndWales = filterEventsOnBunting(s.HolidayData.EnglandAndWales.Events, bunting, divisionName)
	case "scotland":
		filteredHoliday.Scotland = filterEventsOnBunting(s.HolidayData.Scotland.Events, bunting, divisionName)
	case "northern-ireland":
		filteredHoliday.NorthernIreland = filterEventsOnBunting(s.HolidayData.NorthernIreland.Events, bunting, divisionName)
	default:
		http.Error(w, "Invalid division parameter", http.StatusBadRequest)
		return
	}

	writeHttpResponse(w, filteredHoliday)
}

// Handler function for the "/get_year" path
func (s *Server) filterYearHandler(w http.ResponseWriter, r *http.Request) {
	klog.Infoln("Serving req")
	// Parse the "year" query parameter from the request
	yearParam := r.URL.Query().Get("year")

	year, err := strconv.Atoi(yearParam)
	if err != nil {
		http.Error(w, "Invalid year parameter", http.StatusBadRequest)
		return
	}
	filteredHoliday := filterYear(w, r, year, []Division{s.HolidayData.EnglandAndWales, s.HolidayData.Scotland, s.HolidayData.NorthernIreland})

	writeHttpResponse(w, filteredHoliday)

}

// Handler function for the "/get_year_less_event" path
func (s *Server) filterYearWithRestrictedEventHandler(w http.ResponseWriter, r *http.Request) {
	klog.Infoln("Serving req")
	// Parse the "year" query parameter from the request
	yearParam := r.URL.Query().Get("year")

	year, err := strconv.Atoi(yearParam)
	if err != nil {
		http.Error(w, "Invalid year parameter", http.StatusBadRequest)
		return
	}
	filteredHoliday := filterYear(w, r, year, []Division{s.HolidayData.EnglandAndWales, s.HolidayData.Scotland, s.HolidayData.NorthernIreland})

	// Marshal the filteredHoliday into JSON
	jsonData, err := json.Marshal(filteredHoliday)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	var data interface{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusInternalServerError)
		return
	}

	//remove bunting and notes
	removeFields(&data, []string{"bunting", "notes"})
	writeHttpResponse(w, data)
}
