package properties

const (
	BOUNDING_BOX_BYTES = 10
	TILE_DIRECTION_END_BYTES = 1
	GROUNDSPRITE_BYTES = 4
)

type Property interface {
	GetBytes() int
	GetString() string
}





