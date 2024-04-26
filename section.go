package inifile

// A Section represents a single section of an INI file.
//
// The map keys are all lowercase and trimmed of whitespace.
type Section map[string]string

// Set sets the value for a key in a section.
//
// If join is nonzero, the value is appended to any pre-existing value using join as the separator.
//
// Key strings are case-insensitive and ignore leading and trailing whitespace.
func (sect Section) Set(key, value string, join rune) {
	key = Key(key)
	if join != 0 {
		if prev, ok := sect[key]; ok {
			value = prev + string(join) + value
		}
	}
	sect[key] = value
}

// Get looks up a value for a key in a section and returns that value, along with a boolean result similar to a map lookup.
//
// Key strings are case-insensitive and ignore leading and trailing whitespace.
func (sect Section) Get(key string) (value string, ok bool) {
	value, ok = sect[Key(key)]
	return
}

// GetDefault calls Get but returns dflt if the key was not found.
//
// Key strings are case-insensitive and ignore leading and trailing whitespace.
func (sect Section) GetDefault(key, dflt string) string {
	if value, ok := sect.Get(key); ok {
		return value
	}
	return dflt
}
