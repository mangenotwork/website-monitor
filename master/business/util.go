package business

func AlertRuleCode(code int) bool {
	if code == 400 || code == 404 || code >= 500 {
		return true
	}
	return false
}
