package display

type OriginType string

const (
	TopLeft = OriginType("topleft")
	BottomRight = OriginType("bottomright")
	//BottomLeft = OriginType("bottomleft")
)

type Message struct {
	Origin OriginType `json:"origin"`
	Pixels string `json:"pixels"`
	ClearOnMissing bool `json:"clear_on_missing"`
}
