package decision

// Decision matrics detail
type Metrics struct {
	Attr string
	Key1 int32
	Key2 int32
}

// Decision paramter time
type Decision []Metrics

// Helps to find - give value has hight rating to decide
func (d Decision) For(attr string) (ok bool) {
	out := make(map[string]bool)
	for _, m1 := range d {
		out[m1.Attr] = true
		for _, m2 := range d {
			if m1.Attr != m2.Attr {
				if m1.Key1 < m2.Key1 {
					out[m1.Attr] = false
					break
				} else if m1.Key1 == m2.Key1 {
					if m1.Key2 > m2.Key2 {
						out[m1.Attr] = false
						break
					}
				}
			}
		}
	}
	_ok1, _ok2 := out[attr]
	return _ok1 && _ok2
}
