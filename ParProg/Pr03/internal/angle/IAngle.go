package angle

type IAngle interface {
	GetDegree() (float64, error)
	GetRadian() (float64, error)
	Equal(IAngle) bool
}
