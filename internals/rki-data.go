package internals

import (
	"encoding/csv"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//
type RKIData struct {
	Items []RKIDataItem
}

//
type RKIDataItem struct {
	Province string
	Country  string
	Lat      string
	Lng      string
	Date     string
	Cases    string
	Deaths string
}

//
func (r *RKIData) Fetch(url string) error {
	// load the html from the given url
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	//
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	//
	now := time.Now()

	//
	doc.Find("#main table").Each(func(index int, tableHtml *goquery.Selection) {
		//
		rows := tableHtml.Find("tbody tr")

		rows.Each(func(index int, rowHtml *goquery.Selection) {
			// omit the last row in the table
			if index < rows.Length()-1 {
				// create new item
				rkiDataItem := RKIDataItem{
					Country: "Germany",
					Date:    now.Format("1/2/06"),
				}

				//
				rowHtml.Find("td").Each(func(index int, cellHtml *goquery.Selection) {
					//
					switch index {
					case 0:
						//
						lat, lng := provinceToLatLng(cellHtml.Text())
						rkiDataItem.Lat = lat
						rkiDataItem.Lng = lng
						//
						rkiDataItem.Province = provinceI18n(cellHtml.Text())
					case 1:
						//
						casesRaw := cellHtml.Text()
						casesRunes := []rune(casesRaw)
						startIndex := strings.Index(casesRaw, "(")
						endIndex := strings.Index(casesRaw, ")")

						// cases. does the cell include deaths at all?
						if startIndex > 0 {
							rkiDataItem.Cases = string(casesRunes[0:startIndex - 1])
						} else {
							rkiDataItem.Cases = casesRaw
						}

						// deaths. does the cell include deaths at all?
						if startIndex > 0 {
							rkiDataItem.Deaths = string(casesRunes[startIndex + 1:endIndex])
						} else {
							rkiDataItem.Deaths = "0"
						}
					}
				})

				//
				r.Items = append(r.Items, rkiDataItem)
			}
		})
	})

	return nil
}

//
func (r *RKIData) Save(filename string) error {
	//
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	// create the writer
	writer := csv.NewWriter(file)

	// get the csv header
	header := csvHeader()
	// write the header
	err = writer.Write(header)
	if err != nil {
		return err
	}

	//
	total := 0

	//
	for _, item := range r.Items {
		//
		csvString := []string{
			item.Province,
			item.Country,
			item.Lat,
			item.Lng,
			item.Cases,
		}

		cases, err := strconv.Atoi(item.Cases)
		total += cases

		// write the record
		err = writer.Write(csvString)
		if err != nil {
			return err
		}
	}

	totalString := []string{
		"Gesamt",
		"",
		"",
		"",
		strconv.Itoa(total),
	}

	// write the record
	err = writer.Write(totalString)
	if err != nil {
		return err
	}

	// write the buffer to the file
	writer.Flush()
	err = writer.Error()
	if err != nil {
		return err
	}

	return file.Close()
}

// get the header for the csv file
// used in mode fetch-rki-only
func csvHeader() []string {
	//
	now := time.Now()
	//
	return []string{
		"Province",
		"Country",
		"Lat",
		"Lng",
		now.Format("1/2/06"),
	}
}

//
func provinceI18n(province string) string {
	//
	switch province {
	case "Baden-Württemberg":
		return "Baden-Wuerttemberg"
	case "Bayern":
		return "Bavaria"
	case "Berlin":
		return "Berlin"
	case "Brandenburg":
		return "Brandenburg"
	case "Bremen":
		return "Bremen"
	case "Hamburg":
		return "Hamburg"
	case "Hessen":
		return "Hesse"
	case "Mecklenburg-Vorpommern":
		return "Mecklenburg-Vorpommern"
	case "Niedersachsen":
		return "Lower Saxony"
	case "Nordrhein-Westfalen":
		return "North Rhine Westphalia"
	case "Rheinland-Pfalz":
		return "Rhineland-Palatinate"
	case "Saarland":
		return "Saarland"
	case "Sachsen":
		return "Saxony"
	case "Sachsen-Anhalt":
		return "Saxony-Anhalt"
	case "Schleswig-Holstein":
		return "Schleswig-Holstein"
	case "Thüringen":
		return "Thuringia"
	}

	return ""
}

//
func provinceToLatLng(province string) (string, string) {
	//
	switch province {
	case "Baden-Württemberg":
		return "48.537778", "9.041111"
	case "Bayern":
		return "48.7775", "11.431111"
	case "Berlin":
		return "52.52", "13.405"
	case "Brandenburg":
		return "52.361944", "13.008056"
	case "Bremen":
		return "53.075833", "8.8075"
	case "Hamburg":
		return "53.565278", "10.001389"
	case "Hessen":
		return "50.666111", "8.591111"
	case "Mecklenburg-Vorpommern":
		return "53.616667", "12.7"
	case "Niedersachsen":
		return "52.756111", "9.393056"
	case "Nordrhein-Westfalen":
		return "51.466667", "7.55"
	case "Rheinland-Pfalz":
		return "49.913056", "7.45"
	case "Saarland":
		return "49.383056", "6.833056"
	case "Sachsen":
		return "51.026944", "13.358889"
	case "Sachsen-Anhalt":
		return "51.971111", "11.47"
	case "Schleswig-Holstein":
		return "54.47", "9.513889"
	case "Thüringen":
		return "50.861111", "11.051944"
	}

	return "", ""
}
