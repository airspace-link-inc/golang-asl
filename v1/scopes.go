package aslv1

type Scope int8

const (
	AdvisoryRead Scope = iota
	SurfaceTier1
)

func (s Scope) String() string {
	switch s {
	case AdvisoryRead:
		return "airhub-api/advisory.read"
	case SurfaceTier1:
		return "airhub-api/surface.tier1"
	default:
		return "not-found"
	}
}