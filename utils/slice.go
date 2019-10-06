package utils

func UniqueSlice(in []string, transforms ...func(string) string) []string {
	m := map[string]bool{}
	for _, v := range in {
		for _, transform := range transforms {
			v = transform(v)
		}
		m[v] = true
	}
	out := make([]string, 0)
	for v := range m {
		out = append(out, v)
	}
	return out
}
