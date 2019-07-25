package nci_geo_lib

type unit struct {
	to_meter float64
}

var units = map[string]unit{
	"ft":    unit{to_meter: 0.3048},
	"us-ft": unit{to_meter: 1200. / 3937.},
}
