package gcsfuse

import (
	"fmt"
)

func checkRequiredCreateOptions(opts map[string]string, req []string) error {
	notFound := []string{}
	for _, key := range req {
		if _, ok := opts[key]; !ok {
			notFound = append(notFound, key)
		}
	}
	if len(notFound) > 0 {
		return fmt.Errorf("missing required options: %v", notFound)
	}
	return nil
}
