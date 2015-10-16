// Package iniflag is an abstraction around standard go's flag,
// environment variables, and ini configuration reader.
package iniflag

import (
	"flag"
)

// Parse is an equivalent of flag.Parse function but if some
// command argument is missing a value from ini configuration file
// is used.
func Parse(file ...string) error {
	return nil
}

// ParseSet ...
func ParseSet(fset *flag.FlagSet) error {
	return nil
}
