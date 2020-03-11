package internals

// save the values for the csv-header
// of the fetched original csv-file
type CoVDataCSVHeader struct {
	Province string
	Country  string
	Lat      string
	Lng      string
	Dates    []string
}

// check whether the given date are
func (h *CoVDataCSVHeader) HasDate(date string) bool {
	//
	for _, headerDate := range h.Dates {
		//
		if headerDate == date {
			return true
		}
	}

	return false
}

//
func (h *CoVDataCSVHeader) CSVHeader() []string {
	// add the default columns
	columns := []string{
		h.Province,
		h.Country,
		h.Lat,
		h.Lng,
	}

	//
	for _, date := range h.Dates {
		//
		columns = append(columns, date)
	}

	return columns
}
