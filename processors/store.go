package processors

// IStore is a key value permanent storage
// each key/value is saved in a named bucket
// the default bucket name is "default"
type IStore interface {
	Get(string, string) ([]byte, error) // Get(key, bucket)
	Set(string, string, []byte) error   // Set(key, bucket, value)
	Delete(string, string) error        // Delete(key, bucket)
	Has(string, string) (bool, error)   // Has(key, bucket)
}
