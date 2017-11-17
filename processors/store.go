package processors

// IStore is a key value permanent storage
// each key/value is saved in a named bucket
// the default bucket name is "default"
type IStore interface {
	Get(string, string) []byte  // Get(key, bucket)
	Set(string, string, []byte) // Set(key, bucket, value)
	Delete(string, string)      // Delete(key, bucket)
	Has(string, string) bool    // Has(key, bucket)
}
