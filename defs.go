package nci_geo_lib

//存储添加进来的WKT
var defs map[string]*SR

func addDef(name, def string) error {
	if defs == nil {
		defs = make(map[string]*SR)
	}
	proj, err := ParseWKT(def)
	if err != nil {
		return err
	}
	defs[name] = proj
	return nil
}
