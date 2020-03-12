package main

import (
	"fmt"
	"github.com/covid19-timeseries/internals"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

//
const Version = ""
const Build = ""

//
const RkiUrl = "https://www.rki.de/DE/Content/InfAZ/N/Neuartiges_Coronavirus/Fallzahlen.html"
const CoVCsvUrl = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_19-covid-Confirmed.csv"

//
func main() {
	//
	var filename string
	var fetchRKIOnly bool
	//
	app := &cli.App{
		Name:        "covid19-timeseries",
		Version:     fmt.Sprintf("%s-%s", Version, Build),
		Description: "Update the SARS CoV timeseries data with them from the German RKI (Robert Koch Institut)",
		Usage:       "",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "output",
				Usage:       "The filename of the output file.",
				Value:       "./time_series_19-covid-Confirmed.csv",
				Aliases:     []string{"o"},
				Destination: &filename,
			},
			&cli.BoolFlag{
				Name:        "fetch-rki-only",
				Usage:       "Whether to only fetch the data from the rki without merging with existing data.",
				Value:       false,
				Aliases:     []string{"f"},
				Destination: &fetchRKIOnly,
			},
		},
		Action: func(c *cli.Context) error {
			//
			var err error

			// only fetch from rki?
			if fetchRKIOnly {
				//
				err = fetchFromRKI(filename)
				if err != nil {
					log.Printf("Error: %s", err)
				}

				return nil
			}

			// process normally
			err = mergeWithCoVData(filename)
			if err != nil {
				log.Printf("Error: %s", err)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// fetch and merge the existing file and
// output it to a csv file.
func mergeWithCoVData(filename string) error {
	//
	var err error
	//
	p := internals.Parser{
		RkiUrl:    RkiUrl,
		CoVCsvUrl: CoVCsvUrl,
	}

	// fetch current file
	err = p.Fetch()
	if err != nil {
		log.Printf("Error: %s", err)
	}
	originalValuesLen := len(p.Data.Items)
	log.Printf("Done! Found %d existing entries...", originalValuesLen)

	// update with rki data
	err = p.Update()
	if err != nil {
		log.Printf("Error: %s", err)
	}

	// save as csv
	err = p.Save(filename)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	log.Printf("Finished. Created CSV file under '%s'", filename)

	return nil
}

// only fetch the data from the rki and
// output it to a csv file
func fetchFromRKI(filename string) error {
	//
	log.Printf("Fetching rki data from -> %s", RkiUrl)
	//
	rkiData := internals.RKIData{}
	err := rkiData.Fetch(RkiUrl)
	if err != nil {
		return err
	}

	// save as csv
	err = rkiData.Save(filename)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	log.Printf("Finished. Created CSV file under '%s'", filename)

	return nil
}
