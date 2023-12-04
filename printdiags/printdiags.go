package printdiags

import (
	"fmt"

	"github.com/common-fate/clio"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	"github.com/common-fate/sdk/uid"
)

func Print(diags []*accessv1alpha1.Diagnostic, names map[uid.UID]string) bool {
	var haserrors bool
	for _, w := range diags {
		msg := w.Message
		if w.Resource != nil {
			id := uid.FromAPI(w.Resource)
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
