package utils

import (
	"slices"

	"github.com/informatik-mannheim/cmg-ss2025/services/job/ports"
)

// containsStatus checks if a given status is present in the status list.
// It returns true if the status is found, otherwise false.
func ContainsStatus(statusList []ports.JobStatus, status ports.JobStatus) bool {
	return slices.Contains(statusList, status)
}
