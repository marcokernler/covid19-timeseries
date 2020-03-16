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
const CoVCasesCsvUrl = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_19-covid-Confirmed.csv"
const CoVDeathsCsvUrl = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_19-covid-Deaths.csv"

//
func main() {
	//
	var casesFilename string
	var deathsFilename string
	var rkiFilename string
	var fetchRKIOnly bool
	//
	app := &cli.App{
		Name:        "covid19-timeseries",
		Version:     fmt.Sprintf("%s-%s", Version, Build),
		Description: "Update the SARS CoV timeseries data with them from the German RKI (Robert Koch Institut)",
		Usage:       "",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "cases-output",
				Usage:       "The filename of the cases output file.",
				Value:       "./time_series_19-covid-Confirmed.csv",
				// Aliases:     []string{"co"},
				Destination: &casesFilename,
			},
			&cli.StringFlag{
				Name:        "deaths-output",
				Usage:       "The filename of the deaths output file.",
				Value:       "./time_series_19-covid-Deaths.csv",
				// Aliases:     []string{"do"},
				Destination: &deathsFilename,
			},
			&cli.BoolFlag{
				Name:        "fetch-rki-only",
				Usage:       "Whether to only fetch the data from the rki without merging with existing data.",
				Value:       false,
				// Aliases:     []string{"f"},
				Destination: &fetchRKIOnly,
			},
			&cli.StringFlag{
				Name:        "rki-output",
				Usage:       "The filename of the rki file. (only available while in fetch-rki-only mode)",
				Value:       "./rki.csv",
				// Aliases:     []string{"ro"},
				Destination: &rkiFilename,
			},
		},
		Action: func(c *cli.Context) error {
			//
			var err error

			// only fetch from rki?
			if fetchRKIOnly {
				//
				err = fetchFromRKI(rkiFilename)
				if err != nil {
					log.Printf("Error: %s", err)
				}

				return nil
			}

			// process normally
			err = mergeWithCoVData(casesFilename, deathsFilename)
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
func mergeWithCoVData(casesFilename string, deathsFilename string) error {
	//
	var err error
	//
	p := internals.Parser{
		RkiUrl:          RkiUrl,
		CoVCasesCsvUrl:  CoVCasesCsvUrl,
		CoVDeathsCsvUrl: CoVDeathsCsvUrl,
	}

	// fetch current cases file
	err = p.FetchCases()
	if err != nil {
		log.Printf("Error: %s", err)
	}

	// fetch current deaths file
	err = p.FetchDeaths()
	if err != nil {
		log.Printf("Error: %s", err)
	}

	// update with rki data
	err = p.Update()
	if err != nil {
		log.Printf("Error: %s", err)
	}

	// save cases as csv
	err = p.SaveCases(casesFilename)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	// save cases as csv
	err = p.SaveDeaths(deathsFilename)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	log.Printf("Finished. Created CSV files under '%s' and '%s'", casesFilename, deathsFilename)

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
