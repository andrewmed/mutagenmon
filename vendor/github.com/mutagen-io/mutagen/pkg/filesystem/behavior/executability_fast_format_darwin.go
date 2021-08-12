package behavior

import (
	"github.com/mutagen-io/mutagen/pkg/filesystem/behavior/internal/format"
)

// probeExecutabilityPreservationFastByFormat checks if the specified format
// matches well-known executability preservation behavior.
func probeExecutabilityPreservationFastByFormat(f format.Format) (bool, bool) {
	switch f {
	case format.FormatAPFS:
		return true, true
	case format.FormatHFS:
		return true, true
	default:
		return false, false
	}
}
