package printdiags

import (
	"fmt"

	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/eid"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
)

func Print(diags []*accessv1alpha1.Diagnostic, names map[eid.EID]string) bool {
	var haserrors bool
	for _, w := range diags {
		msg := w.Message
		if w.Resource != nil {
			id := eid.FromAPI(w.Resource)
			msg = fmt.Sprintf("%s: %s", id, msg)
		}

		switch w.Level {
		case accessv1alpha1.DiagnosticLevel_DIAGNOSTIC_LEVEL_WARNING:
			haserrors = true
			clio.Warn(msg)
		case accessv1alpha1.DiagnosticLevel_DIAGNOSTIC_LEVEL_ERROR:
			haserrors = true
			clio.Error(msg)
		default:
			clio.Info(msg)
		}
	}
	return haserrors
}
