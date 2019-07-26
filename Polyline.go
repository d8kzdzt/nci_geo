package nci_geo_lib

//定义折线
type PolyLine []Point

//表示一个极大值
const MaxValue float64 = 99999999999999999999999999999999999999999999999999999999999999

//返回折线总长度
func (self *PolyLine) TotalDistance() float64 {
	_len := len(*self)
	if _len == 0 {
		return 0
	}
	return self.Distance(_len - 1)
}

//返回第idx个点(idx从0开始)距起点的距离之和
func (self *PolyLine) Distance(idx int) float64 {
	dist := 0.0
	if idx > len(*self)+1 {
		panic("idx超出折线的长度")
	}
	for i := 0; i < idx; i++ {
		dist += (*self)[i].Distance(&(*self)[i+1])
	}
	return dist
}

//获取【折线总长,投影点前面一个点下标(从0开始),指定点到折线起点距离,指定点到投影点距离，投影点,指定点在线段右侧?】
func (self *PolyLine) QueryPointAndDistance(pt *Point) (float64, int, float64, float64, Point, bool) {
	if len(*self) <= 1 {
		panic("polyLine至少包含2个点")
	}
	var (
		minDistPoint  Point   //投影点即最短距离点
		distFromBegin float64 //投影点距起点距离
		lastPxIdx     int     //投影点前面一个点下标(从0开始)
		isOnRight     bool    //指定点在折线的右侧？
	)
	endPoint := (*self)[len(*self)-1]
	minDist := MaxValue
	//遍历每个线段，求最小距离所在的线段
	for i := 0; i < len(*self)-1; i++ {
		//点到线段的距离
		dist := DistanceFromSegment(pt, &(*self)[i], &(*self)[i+1])
		if dist < minDist {
			//最短距离被刷新
			minDist = dist
			lastPxIdx = i
			minDistPoint = ClosestPoint(pt, &(*self)[i], &(*self)[i+1])
		}
	}
	plLen := self.TotalDistance()
	isOnRight = IsOnRight(pt, &(*self)[lastPxIdx], &(*self)[lastPxIdx+1])
	distFromBegin = self.Distance(lastPxIdx) + minDistPoint.Distance(&(*self)[lastPxIdx])
	//投影点时最后一个点时，上一个点下标置为最后一个点
	if minDistPoint.Equals(&endPoint){
		lastPxIdx = len(*self) -1
	}
	return plLen, lastPxIdx, distFromBegin, minDist, minDistPoint, isOnRight
}
