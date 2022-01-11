/*
 This package helps defining structures that are shared across backends
 so as to avoid import cyles
*/
package types

// Metadata structures the information we expect from the objects
type Metadata struct {
	Version string
	Source  string
}
