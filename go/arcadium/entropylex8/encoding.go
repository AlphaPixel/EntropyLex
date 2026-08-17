package entropylex8

type (
	Encoding struct {
		dict map[int]string
	}
)

func NewEncoding(dict map[int]string) *Encoding {
	return &Encoding{dict: dict}
}
