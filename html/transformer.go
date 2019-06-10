package html

type Transformer interface {
	Apply() map[string][]byte
}
