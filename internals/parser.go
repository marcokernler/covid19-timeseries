package internals

import (
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"os"
)

//
type Parser struct {
	RkiUrl     string
	CoVCasesCsvUrl  string
	CoVDeathsCsvUrl string
	HeaderData CoVDataCSVHeader
	Data       CoVData
}

// fetch the current cases file from github and
// parse the data into our data-structure
func (p *Parser) FetchCases() error {
	//
	log.Printf("Fetching existing cases data from -> %s", p.CoVCasesCsvUrl)

	//
	res, err := http.Get(p.CoVCasesCsvUrl)
	if err != nil {
		return err
	}

	//
	reader := csv.NewReader(res.Body)
	covDataItems := []CoVDataItem{}

	// read the header of the csv
	header, err := reader.Read()
	if err != nil {
		return err
	}

	// persist the static values in the header
	p.HeaderData.Province = header[0]
	p.HeaderData.Country = header[1]
	p.HeaderData.Lat = header[2]
	p.HeaderData.Lng = header[3]

	// loop over the records
	for i := 0; ; i++ {
		//
		record, err := reader.Read()
		covDataItem := CoVDataItem{}

		// check if we're on the end
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		// loop over the columns
		for col := 0; col < len(record); col++ {
			// the first 4 cols are fixed data
			if col < 4 {
				switch col {
				case 0:
					covDataItem.Province = record[col]
				case 1:
					covDataItem.Country = record[col]

					// log.Printf("item: %d - %s", i, covDataItem.Country)
				case 2:
					covDataItem.Lat = record[col]
				case 3:
					covDataItem.Lng = record[col]
				}
			} else {
				// while on the first row, collect the date values
				// from the header of the csv file.
				if i == 0 {
					p.HeaderData.Dates = append(p.HeaderData.Dates, header[col])
				}

				//
				value := CoVDataItemValue{
					Date:  header[col],
					Cases: record[col],
				}

				covDataItem.Values = append(covDataItem.Values, value)
			}
		}

		covDataItems = append(covDataItems, covDataItem)
	}

	//
	p.Data.Items = covDataItems

	return nil
}

// fetch the current deaths file from github and
// parse the data into our data-structure. Currently
// fetching the cases needs to be done first, since
// FetchCases will create the base data structure with
// the dates for the csv file.
func (p *Parser) FetchDeaths() error {
	//
	log.Printf("Fetching existing deaths data from -> %s", p.CoVDeathsCsvUrl)

	//
	res, err := http.Get(p.CoVDeathsCsvUrl)
	if err != nil {
		return err
	}

	//
	reader := csv.NewReader(res.Body)

	// read the header of the csv
	_, err = reader.Read()
	if err != nil {
		return err
	}

	// loop over the records
	for i := 0; ; i++ {
		// read the current row
		record, err := reader.Read()

		// check if we're on the end
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		// get item from existing data
		covDataItem := &p.Data.Items[i]

		// loop over the columns
		for col := 0; col < len(record); col++ {
			// only
			if col >= 4 {
				// get the value and update the death count
				itemValue := &covDataItem.Values[col - 4]
				itemValue.Deaths = record[col]
			}
		}
	}

	return nil
}

// update the current data model with the data
// from the rki
func (p *Parser) Update() error {
	//
	log.Printf("Fetching data from -> %s", p.RkiUrl)
	//
	rkiData := RKIData{}
	err := rkiData.Fetch(p.RkiUrl)
	if err != nil {
		return err
	}

	//
	for _, rkiDataItem := range rkiData.Items {
		// if the header doesn't contain the
		// new date, add to the header
		if !p.HeaderData.HasDate(rkiDataItem.Date) {
			//
			p.HeaderData.Dates = append(p.HeaderData.Dates, rkiDataItem.Date)
		}

		// try getting existing data item
		covDataItem := p.Data.getByProvince(rkiDataItem.Province)

		// if no entry exists for the province
		// create a new record
		if covDataItem.Province == "" {
			//
			covDataItem = CoVDataItem{
				Province: rkiDataItem.Province,
				Country:  rkiDataItem.Country,
				Lat:      rkiDataItem.Lat,
				Lng:      rkiDataItem.Lng,
			}

			// fill up missing dates
			for i := 0; i < len(p.HeaderData.Dates)-1; i++ {
				//
				date := p.HeaderData.Dates[i]
				//
				covDataItem.Values = append(covDataItem.Values, CoVDataItemValue{
					Date:  date,
					Cases: "0",
					Deaths: "0",
				})
			}

			// add the new value
			covDataItem.Values = append(covDataItem.Values, CoVDataItemValue{
				Date:  rkiDataItem.Date,
				Cases: rkiDataItem.Cases,
				Deaths: rkiDataItem.Deaths,
			})

			// add the item to the list
			p.Data.Items = append(p.Data.Items, covDataItem)
		} else {
			// add the rki-value to the existing item
			covDataItem.Values = append(covDataItem.Values, CoVDataItemValue{
				Date:  rkiDataItem.Date,
				Cases: rkiDataItem.Cases,
				Deaths: rkiDataItem.Deaths,
			})
		}
	}

	return nil
}

// save the current data model to the a file
// with the given filename for the case count
func (p *Parser) SaveCases(filename string) error {
	//
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	// create the writer
	writer := csv.NewWriter(file)

	// get the csv header
	header := p.HeaderData.CSVHeader()
	// write the header
	err = writer.Write(header)
	if err != nil {
		return err
	}

	//
	for _, item := range p.Data.Items {
		//
		csvString := []string{
			item.Province,
			item.Country,
			item.Lat,
			item.Lng,
		}

		// loop over the values of the item
		for _, values := range item.Values {
			// append as column
			csvString = append(csvString, values.Cases)
		}

		// write the record
		err = writer.Write(csvString)
		if err != nil {
			return err
		}
	}

	// write the buffer to the file
	writer.Flush()
	err = writer.Error()
	if err != nil {
		return err
	}

	return file.Close()
}

// save the current data model to the a file
// with the given filename for the death count
func (p *Parser) SaveDeaths(filename string) error {
	//
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	// create the writer
	writer := csv.NewWriter(file)

	// get the csv header
	header := p.HeaderData.CSVHeader()
	// write the header
	err = writer.Write(header)
	if err != nil {
		return err
	}

	//
	for _, item := range p.Data.Items {
		//
		csvString := []string{
			item.Province,
			item.Country,
			item.Lat,
			item.Lng,
		}

		// loop over the values of the item
		for _, values := range item.Values {
			// append as column
			csvString = append(csvString, values.Deaths)
		}

		// write the record
		err = writer.Write(csvString)
		if err != nil {
			return err
		}
	}

	// write the buffer to the file
	writer.Flush()
	err = writer.Error()
	if err != nil {
		return err
	}

	return file.Close()
}
