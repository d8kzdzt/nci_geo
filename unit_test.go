package nci_geo_lib

import (
	"fmt"
	"testing"
)

//测试投影点，距起点距离等
func TestPolyLineAndDistance(t *testing.T) {
	//生成一个100个点组成的折线
	var polyLine PolyLine
	polyLine = append(polyLine, Point{1, 1})
	polyLine = append(polyLine, Point{2, 1})
	polyLine = append(polyLine, Point{3, 1})
	polyLine = append(polyLine, Point{4, 1})
	polyLine = append(polyLine, Point{5, 1})

	//目标点
	p := Point{6, -3}
	//测试
	plLen, projectIdx, distFromBegin, minDist, minDistPoint, onRight := polyLine.QueryPointAndDistance(&p)
	fmt.Printf("折线总长:%.3f\n", plLen)
	fmt.Printf("投影位置:%d\n", projectIdx)
	fmt.Printf("距起点:%.3f\n", distFromBegin)
	fmt.Printf("距投影点:%.3f\n", minDist)
	fmt.Printf("投影点:[X:%.3f,Y:%.3f]\n", minDistPoint.X, minDistPoint.Y)
	fmt.Printf("在右侧?:%v\n", onRight)
}

//gps坐标转WKT投影坐标
func TestCoordinateConvert(t *testing.T) {
	x, y := 119.086825, 31.235314
	var (
		a, b           float64
		err            error
		cs4326, cs2437 *SR
		trans          Transformer
	)
	//WKT标准的两个坐标系描述字符串
	wkt4326 := "GEOGCS[\"WGS 84\",DATUM[\"WGS_1984\",SPHEROID[\"WGS 84\",6378137,298.257223563,AUTHORITY[\"EPSG\",\"7030\"]],AUTHORITY[\"EPSG\",\"6326\"]],PRIMEM[\"Greenwich\",0,AUTHORITY[\"EPSG\",\"8901\"]],UNIT[\"degree\",0.01745329251994328,AUTHORITY[\"EPSG\",\"9122\"]],AUTHORITY[\"EPSG\",\"4326\"]]"
	wkt2437 := "PROJCS[\"Beijing 1954 / 3-degree Gauss-Kruger CM 120E\",GEOGCS[\"Beijing 1954\",DATUM[\"Beijing_1954\",SPHEROID[\"Krassowsky 1940\",6378245,298.3,AUTHORITY[\"EPSG\",\"7024\"]],TOWGS84[15.8,-154.4,-82.3,0,0,0,0],AUTHORITY[\"EPSG\",\"6214\"]],PRIMEM[\"Greenwich\",0,AUTHORITY[\"EPSG\",\"8901\"]],UNIT[\"degree\",0.0174532925199433,AUTHORITY[\"EPSG\",\"9122\"]],AUTHORITY[\"EPSG\",\"4214\"]],PROJECTION[\"Transverse_Mercator\"],PARAMETER[\"latitude_of_origin\",0],PARAMETER[\"central_meridian\",120],PARAMETER[\"scale_factor\",1],PARAMETER[\"false_easting\",500000],PARAMETER[\"false_northing\",0],UNIT[\"metre\",1,AUTHORITY[\"EPSG\",\"9001\"]],AUTHORITY[\"EPSG\",\"2437\"]]"
	//解析两个坐标系到内存中
	if cs4326, err = ParseWKT(wkt4326); err != nil {
		t.Fatalf("wkt4326解析错误:%v", err)
	}
	if cs2437, err = ParseWKT(wkt2437); err != nil {
		t.Fatalf("wkt2437解析错误:%v", err)
	}
	//trans表示一个函数，此函数用来将cs4326下的坐标转为cs2437下的坐标
	if trans, err = cs4326.CoordinateTransformNew(cs2437); err != nil {
		t.Fatalf("wkt4326转wkt2437错误:%v", err)
	}
	if a, b, err = trans(x, y); err != nil {
		t.Fatalf("坐标转换失败")
	}
	fmt.Printf("投影坐标:%.5f\t%.5f\n", a, b)
}

func BenchmarkCoordinateConvert(bm *testing.B) {
	x, y := 119.086825, 31.235314
	var (
		err            error
		cs4326, cs2437 *SR
		trans          Transformer
	)
	//WKT标准的两个坐标系描述字符串
	wkt4326 := "GEOGCS[\"WGS 84\",DATUM[\"WGS_1984\",SPHEROID[\"WGS 84\",6378137,298.257223563,AUTHORITY[\"EPSG\",\"7030\"]],AUTHORITY[\"EPSG\",\"6326\"]],PRIMEM[\"Greenwich\",0,AUTHORITY[\"EPSG\",\"8901\"]],UNIT[\"degree\",0.01745329251994328,AUTHORITY[\"EPSG\",\"9122\"]],AUTHORITY[\"EPSG\",\"4326\"]]"
	wkt2437 := "PROJCS[\"Beijing 1954 / 3-degree Gauss-Kruger CM 120E\",GEOGCS[\"Beijing 1954\",DATUM[\"Beijing_1954\",SPHEROID[\"Krassowsky 1940\",6378245,298.3,AUTHORITY[\"EPSG\",\"7024\"]],TOWGS84[15.8,-154.4,-82.3,0,0,0,0],AUTHORITY[\"EPSG\",\"6214\"]],PRIMEM[\"Greenwich\",0,AUTHORITY[\"EPSG\",\"8901\"]],UNIT[\"degree\",0.0174532925199433,AUTHORITY[\"EPSG\",\"9122\"]],AUTHORITY[\"EPSG\",\"4214\"]],PROJECTION[\"Transverse_Mercator\"],PARAMETER[\"latitude_of_origin\",0],PARAMETER[\"central_meridian\",120],PARAMETER[\"scale_factor\",1],PARAMETER[\"false_easting\",500000],PARAMETER[\"false_northing\",0],UNIT[\"metre\",1,AUTHORITY[\"EPSG\",\"9001\"]],AUTHORITY[\"EPSG\",\"2437\"]]"
	//解析两个坐标系到内存中
	if cs4326, err = ParseWKT(wkt4326); err == nil {
		if cs2437, err = ParseWKT(wkt2437); err == nil {
			if trans, err = cs4326.CoordinateTransform(cs2437); err == nil {
				for i := 0; i < bm.N; i++ {
					trans(x, y)
					return
				}
			}
		}
	}
	bm.Fatalf("性能测试错误,%v", err)
}

//优化后的版本，将函数提前查找好，省去每次调用trans的时候再去defs中找转换函数
func BenchmarkCoordinateConvertNew(bm *testing.B) {
	x, y := 119.086825, 31.235314
	var (
		err            error
		cs4326, cs2437 *SR
		trans          Transformer
	)
	//WKT标准的两个坐标系描述字符串
	wkt4326 := "GEOGCS[\"WGS 84\",DATUM[\"WGS_1984\",SPHEROID[\"WGS 84\",6378137,298.257223563,AUTHORITY[\"EPSG\",\"7030\"]],AUTHORITY[\"EPSG\",\"6326\"]],PRIMEM[\"Greenwich\",0,AUTHORITY[\"EPSG\",\"8901\"]],UNIT[\"degree\",0.01745329251994328,AUTHORITY[\"EPSG\",\"9122\"]],AUTHORITY[\"EPSG\",\"4326\"]]"
	wkt2437 := "PROJCS[\"Beijing 1954 / 3-degree Gauss-Kruger CM 120E\",GEOGCS[\"Beijing 1954\",DATUM[\"Beijing_1954\",SPHEROID[\"Krassowsky 1940\",6378245,298.3,AUTHORITY[\"EPSG\",\"7024\"]],TOWGS84[15.8,-154.4,-82.3,0,0,0,0],AUTHORITY[\"EPSG\",\"6214\"]],PRIMEM[\"Greenwich\",0,AUTHORITY[\"EPSG\",\"8901\"]],UNIT[\"degree\",0.0174532925199433,AUTHORITY[\"EPSG\",\"9122\"]],AUTHORITY[\"EPSG\",\"4214\"]],PROJECTION[\"transverse_mercator\"],PARAMETER[\"latitude_of_origin\",0],PARAMETER[\"central_meridian\",120],PARAMETER[\"scale_factor\",1],PARAMETER[\"false_easting\",500000],PARAMETER[\"false_northing\",0],UNIT[\"metre\",1,AUTHORITY[\"EPSG\",\"9001\"]],AUTHORITY[\"EPSG\",\"2437\"]]"
	//解析两个坐标系到内存中
	if cs4326, err = ParseWKT(wkt4326); err == nil {
		if cs2437, err = ParseWKT(wkt2437); err == nil {
			if trans, err = cs4326.CoordinateTransformNew(cs2437); err == nil {
				for i := 0; i < bm.N; i++ {
					trans(x, y)
				}
				return
			}
		}
	}
	bm.Fatalf("性能测试错误,%v", err)
}


//坐标转换,
func TestCoordinateConvertWrapper(t *testing.T) {
	longitude, latitude := 119.086825, 31.235314
	//生成一个坐标转换对象，内部自动初始化，自动解析了gps到平面坐标的转换逻辑
	//调用者要保存conv对象，因为构造conv的消耗是巨大的，其内部包含了解析协议的细节
	//conv对象可重复使用(调conv.Trans)
	conv := new(CoordinateConverter)
	//将longitude, latitude转为a,b表示的平面坐标,失败返回error对象
	a,b,_ := conv.Trans(longitude, latitude)
	fmt.Printf("%.5f\t%.5f\n",a,b)
}
