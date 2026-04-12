package export

import "fmt"

// ParseFormat converts a raw string into a Format constant.
// It returns an error if the value is not recognised.
func ParseFormat(s string) (Format, error) {
	switch Format(s) {
	case FormatCSV, FormatMarkdown:
		return Format(s), nil
	default:
		return "", fmt.Errorf("export: unrecognised format %q; valid values are \"csv\", \"markdown\"", s)
	}
}

// String implements fmt.Stringer.
func (f Format) String() string {
	return string(f)
}

// KnownFormats returns all supported Format values in a stable order.
func KnownFormats() []Format {
	return []Format{FormatCSV, FormatMarkdown}
}
