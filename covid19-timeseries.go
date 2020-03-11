package main

import (
	"fmt"
	"github.com/covid19-timeseries/internals"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var VERSION = "1.0.0"
var BUILD = "0"

func main() {
	//
	var filename string
	//
	app := &cli.App{
		Name:        "covid19-timeseries",
		Version:     fmt.Sprintf("%s-%s", VERSION, BUILD),
		Description: "Update the SARS CoV timeseries data with them from the German RKI (Robert Koch Institut)",
		Usage:       "",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Usage:   "The filename of the output file.",
				Value:   "./time_series_19-covid-Confirmed.csv",
				Aliases: []string{"o"},
				Destination: &filename,
			},
		},
		Action: func(c *cli.Context) error {
			//
			var err error
			//
			p := internals.Parser{
				RkiUrl:    "https://www.rki.de/DE/Content/InfAZ/N/Neuartiges_Coronavirus/Fallzahlen.html",
				CoVCsvUrl: "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_19-covid-Confirmed.csv",
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
			// updatedValuesLen := len(p.Data.Items) - originalValuesLen
			// if updatedValuesLen != originalValuesLen {
			// 		log.Printf("Done! Updated with %d new values...", updatedValuesLen)
			// } else {
			// 		log.Printf("Done! Updated values...")
			// }

			// save as csv
			err = p.Save(filename)
			if err != nil {
				log.Printf("Error: %s", err)
			}

			log.Printf("Finished. Updated CSV file under '%s'", filename)

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
