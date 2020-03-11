# SARS Covid 19 Time Series Updater

This tool will update the timeseries of the [data repository for the 2019 Novel Coronavirus Visual Dashboard](https://github.com/CSSEGISandData/COVID-19).
 
The cli-tool will download the existing [timeseries csv file](https://github.com/CSSEGISandData/COVID-19/blob/master/csse_covid_19_data/csse_covid_19_time_series/time_series_19-covid-Confirmed.csv) and extend/update it with the values from the German RKI (Robert Koch Institute) and save it to your disk.

## Usage

Download the binary for your OS from [here](https://github.com/marcokernler/covid19-timeseries/releases).

```bash
# print help
$ covid19-timeseries help

# update the csv-file an save it under it's original filename
$ covid19-timeseries

# update the csv-file an save it under test.csv
$ covid19-timeseries --output test.csv

# only fetch data from rki and save under rki.csv
$ covid19-timeseries --fetch-rki-only --output rki.csv
```

## Build from Source
```bash
# clone the repo
$ git clone https://github.com/marcokernler/covid19-timeseries
$ cd covid19-timeseries

# build for all platforms
$ make build_all
```

## License
Copyright 2020 Marco Kernler

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
