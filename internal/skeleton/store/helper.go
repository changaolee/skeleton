package store

const defaultLimitValue = 20

func defaultLimit(limit int) int {
	if limit == 0 {
		limit = defaultLimitValue
	}
	return limit
}
