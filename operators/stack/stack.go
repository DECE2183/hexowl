package stack

func Push[Slice ~[]E, E any](dest Slice, elm E) Slice {
	return append(dest, elm)
}

func Pop[Slice ~[]E, E any](src Slice) (Slice, E) {
	l := len(src)
	return src[:l-1], src[l-1]
}
