package utils

type Variables map[string]any

func MergeVariables(list ...Variables) Variables {
	vars := Variables{}
	for _, src := range list {
		if src == nil {
			continue
		}
		for k, v := range src {
			vars[k] = v
		}
	}
	return vars
}
