// Package validate provides deep validation of cron expressions beyond simple
// parsing. While the parser package confirms structural correctness, this
// package checks that each field's values fall within their allowed ranges,
// that range bounds are ordered correctly, and that step values are positive.
//
// Usage:
//
//	res := validate.Check("0 25 * * *")
//	if !res.Valid {
//	    for _, e := range res.Errors {
//	        fmt.Println(e) // invalid hour field "25": 25 out of allowed range 0-23
//	    }
//	}
//
// Supported field syntax:
//
//	*          - wildcard (any value)
//	5          - literal value
//	1,2,3      - comma-separated list
//	1-5        - inclusive range
//	*/5        - step over wildcard
//	1-10/2     - step over range
//
// Field positions and their allowed ranges:
//
//	Position  Field         Range
//	0         minute        0-59
//	1         hour          0-23
//	2         day-of-month  1-31
//	3         month         1-12
//	4         day-of-week   0-6  (0 = Sunday)
package validate
