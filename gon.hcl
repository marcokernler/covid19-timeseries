

source = [
    "./dist/covid19-timeseries-macos_darwin_amd64/covid19-timeseries"
]
bundle_id = "de.marcokernler.covid19-timeseries"

apple_id {
    username = "office@denkfabrik-neuemedien.de"
    password = "@env:AC_PASSWORD"
    provider = "denkfabrik-neueMedien"
}

sign {
    application_identity = "780B1411B559D93F6851EFBC12679B862E233A93"
}

dmg {
    output_path = "./dist/covid19-timeseries-macos_darwin.dmg"
    volume_name = "covid19-timeseries"
}