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
	return self.Distance(_len)
}

//返回折线从起点到指定下标Idx的总长度
func (self *PolyLine) Distance(idx int) float64 {
	dist := 0.0
	if idx > len(*self) {
		panic("idx超出折线的长度")
	}
	for i := 1; i < idx; i++ {
		dist += (*self)[i-1].Distance(&(*self)[i])
	}
	return dist
}

//获取【折线总长,投影点所在线段下标,指定点到折线起点距离,指定点到投影点距离，投影点,指定点在线段右侧?】
func (self *PolyLine) QueryPointAndDistance(pt *Point) (float64, int, float64, float64, Point, bool) {
	if len(*self) <= 1 {
		panic("polyLine至少包含2个点")
	}
	var (
		minDistPoint  Point
		distFromBegin float64
		projectIdx    int
		isOnRight     bool
	)
	minDist := MaxValue
	//遍历每个线段，求最小距离所在的线段
	for i := 1; i < len(*self); i++ {
		//点到线段的距离
		dist := DistanceFromSegment(pt, &(*self)[i-1], &(*self)[i])
		if dist < minDist {
			//最短距离被刷新
			minDist = dist
			projectIdx = i - 1
			//求投影点(最近点)
			minDistPoint = ClosestPoint(pt, &(*self)[i-1], &(*self)[i])
		}
	}
	plLen := self.TotalDistance()
	if projectIdx == 0 {
		distFromBegin = minDist
		isOnRight = IsOnRight(pt, &(*self)[0], &(*self)[1])
	} else {
		distFromBegin = self.Distance(projectIdx) + minDistPoint.Distance(&(*self)[projectIdx-1])
		isOnRight = IsOnRight(pt, &(*self)[projectIdx-1], &(*self)[projectIdx])
	}
	return plLen, projectIdx, distFromBegin, minDist, minDistPoint, isOnRight
}
