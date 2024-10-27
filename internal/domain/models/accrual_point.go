package models

import "fmt"

type AccrualPoint int

const (
	OneHundredthPoint AccrualPoint = 1
	OnePoint                       = 100 * OneHundredthPoint
)

func (a AccrualPoint) String() string {
	points := int(a) / 100
	hundredth := int(a) % 100

	tenth := hundredth / 10
	hundredth = hundredth % 10
	switch {
	case hundredth > 0:
		return fmt.Sprintf("%d.%d%d", points, tenth, hundredth)
	case tenth > 0:
		return fmt.Sprintf("%d.%d", points, tenth)
	default:
		return fmt.Sprintf("%d", points)
	}
}

func (a AccrualPoint) Float64() float64 {
	return float64(a) / 100.0
}

func (a AccrualPoint) MarshalJSON() ([]byte, error) {
	s := a.String()
	return []byte(s), nil
}

func AccrualPointFromFloat64(num float64) AccrualPoint {
	return AccrualPoint(num * 100)
}
