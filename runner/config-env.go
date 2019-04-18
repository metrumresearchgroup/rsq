package runner

import (
	"fmt"
	"strings"
)

// sysEnvVars contains the default environment variables usually from
// os.Environ()
func configureEnv(sysEnvVars []string, rs RSettings) []string {
	envMap := make(map[string]string)
	envVars := []string{}
	envOrder := []string{}

	for k, v := range rs.EnvVars {
		_, exists := envMap[k]
		if !exists {
			envMap[k] = v
			envOrder = append(envOrder, k)
		}
	}
	// system env vars generally
	for _, ev := range sysEnvVars {
		evs := strings.SplitN(ev, "=", 2)
		if len(evs) > 1 && evs[1] != "" {
			// if exists would be set from the user hence should not accept the system env
			_, exists := envMap[evs[0]]
			if !exists {
				envMap[evs[0]] = evs[1]
				envOrder = append(envOrder, evs[0])
			}
		}
	}

	for _, ev := range envOrder {
		val := envMap[ev]
		envVars = append(envVars, fmt.Sprintf("%s=%s", ev, val))
	}

	return envVars
}
