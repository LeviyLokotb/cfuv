package geo

type IDistanceCalc interface {
	DistanceBetween(a, b Coordinate) (float64, error)
}
