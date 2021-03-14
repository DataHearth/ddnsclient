package utils

func AggregateSubdomains(subdomains []string, domain string) []string {
	agdSub := make([]string, len(subdomains))
	for _, sd := range subdomains {
		agdSub = append(agdSub, sd+"."+domain)
	}

	return agdSub
}
