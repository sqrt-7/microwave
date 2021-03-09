package tools

import (
	"fmt"
	"os"
	"strings"
)

func EnvLookup(names ...string) (map[string]string, error) {
	m := make(map[string]string)
	var missing []string

	for _, v := range names {
		item, present := os.LookupEnv(v)
		if !present {
			missing = append(missing, v)
			continue
		}
		m[v] = item
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing environment variable(s) [%s]", strings.Join(missing, ", "))
	}

	return m, nil
}
