package cout

type FontModifier func(s string) string

func Multi(fm ...FontModifier) FontModifier {
	return func(s string) string {
		for _, f := range fm {
			s = f(s)
		}
		return s
	}
}
