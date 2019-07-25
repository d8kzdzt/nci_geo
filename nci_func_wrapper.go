package nci_geo_lib

var trans4326To2437,trans2437To4326 Transformer

//将x,y表示的经纬度转为WKT2437的平面坐标
func WKT4326ToWKT2437(x,y float64)(float64,float64,error){
	return trans4326To2437(x,y)
}

//将x,y表示的平面坐标转为WKT4326的经纬度
func WKT2437ToWKT4326(x,y float64)(float64,float64,error){
	return trans2437To4326(x,y)
}

func init(){
	var err error
	wkt4326 := "GEOGCS[\"WGS 84\",DATUM[\"WGS_1984\",SPHEROID[\"WGS 84\",6378137,298.257223563,AUTHORITY[\"EPSG\",\"7030\"]],AUTHORITY[\"EPSG\",\"6326\"]],PRIMEM[\"Greenwich\",0,AUTHORITY[\"EPSG\",\"8901\"]],UNIT[\"degree\",0.01745329251994328,AUTHORITY[\"EPSG\",\"9122\"]],AUTHORITY[\"EPSG\",\"4326\"]]"
	wkt2437 := "PROJCS[\"Beijing 1954 / 3-degree Gauss-Kruger CM 120E\",GEOGCS[\"Beijing 1954\",DATUM[\"Beijing_1954\",SPHEROID[\"Krassowsky 1940\",6378245,298.3,AUTHORITY[\"EPSG\",\"7024\"]],TOWGS84[15.8,-154.4,-82.3,0,0,0,0],AUTHORITY[\"EPSG\",\"6214\"]],PRIMEM[\"Greenwich\",0,AUTHORITY[\"EPSG\",\"8901\"]],UNIT[\"degree\",0.0174532925199433,AUTHORITY[\"EPSG\",\"9122\"]],AUTHORITY[\"EPSG\",\"4214\"]],PROJECTION[\"Transverse_Mercator\"],PARAMETER[\"latitude_of_origin\",0],PARAMETER[\"central_meridian\",120],PARAMETER[\"scale_factor\",1],PARAMETER[\"false_easting\",500000],PARAMETER[\"false_northing\",0],UNIT[\"metre\",1,AUTHORITY[\"EPSG\",\"9001\"]],AUTHORITY[\"EPSG\",\"2437\"]]"
	var cs4326,cs2437 *SR
	if cs4326, err = ParseWKT(wkt4326); err == nil {
		if cs2437, err = ParseWKT(wkt2437); err == nil {
			if trans4326To2437, err = cs4326.CoordinateTransform(cs2437); err == nil {
				if trans2437To4326, err = cs2437.CoordinateTransform(cs4326); err == nil {
				}
			}
		}
	}
	if err !=nil{
		panic("坐标转换初始化失败")
	}
}