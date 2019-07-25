package nci_geo_lib

func init() {
	err := addDef("EPSG:4326", "+title=WGS 84 (long/lat) +proj=longlat +ellps=WGS84 +datum=WGS84 +units=degrees")
	if err != nil {
		panic(err)
	}
	defs["WGS84"] = defs["EPSG:4326"]
}
