package postgres

import "fmt"

// PlaceHolder try to build a place holder and arguments for pq library (postgres)
func PlaceHolder(start int, params ...interface{}) ([]string, []interface{}) {
	var res []string
	var p []interface{}

	for i := range params {
		res = append(res, fmt.Sprintf("$%d", start+i))
		p = append(p, params[i])
	}

	return res, p
}
