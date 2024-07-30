package subdomain

import "time"

// SubIsPending check if the current subdomain is waiting the DNS propagation.
func (sb *subdomain) SubIsPending(sbs PendingSubdomains) bool {
	for _, sub := range sbs {
		if sb == sub {
			return true
		}
	}

	return false
}

// FindSubdomain returns a subdomain found in the pending map of subdomain.
// If not found, it returns nil.
func (sb *subdomain) FindSubdomain(sbs PendingSubdomains) Subdomain {
	for _, sub := range sbs {
		if sub == sb {
			return sb
		}
	}

	return nil
}

// CheckPendingSubdomains check if any pending subdomains are waiting to be restored.
// If so, it/they will be returned as a slice.
// If not, it returns nil.
func CheckPendingSubdomains(sbs PendingSubdomains, now time.Time) PendingSubdomains {
	delSbs := make(PendingSubdomains)
	for t, sb := range sbs {
		if t.Add(5 * time.Minute).Before(now) {
			delSbs[t] = sb
		}
	}

	if len(delSbs) < 1 {
		return nil
	}

	return delSbs
}

func DeletePendingSubdomains(delSbs PendingSubdomains, pending PendingSubdomains) PendingSubdomains {
	for t := range delSbs {
		delete(pending, t)
	}

	return pending
}
