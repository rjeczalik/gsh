package sshfile

func nonempty(s ...string) (t []string) {
	for _, s := range s {
		if s != "" {
			t = append(t, s)
		}
	}
	return t
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}
