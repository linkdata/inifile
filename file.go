package inifile

// A File represents a parsed INI file.
// Values found outside of a named section are in the "" section.
//
// The map keys are all lowercase and trimmed of whitespace.
type File map[string]Section

// Section returns a named Section. A Section will be created if one does not already exist for the given name.
//
// Section and key strings are case-insensitive and ignore leading and trailing whitespace.
func (inif File) Section(name string) (section Section) {
	name = Key(name)
	if section = inif[name]; section == nil {
		section = make(Section)
		inif[name] = section
	}
	return
}

// Set sets the value for a key in a section.
//
// The section is created if needed.
// The value will be whitespace trimmed before being stored.
// If join is nonzero, the value is appended to any pre-existing value using join as the separator.
//
// Section and key strings are case-insensitive and ignore leading and trailing whitespace.
func (inif File) Set(section, key, value string, join rune) {
	inif.Section(section).Set(key, value, join)
}

// Get looks up a value for a key in a section and returns that value, along with a boolean result similar to a map lookup.
//
// Section and key strings are case-insensitive and ignore leading and trailing whitespace.
func (inif File) Get(section, key string) (value string, ok bool) {
	return inif[Key(section)].Get(key)
}

// GetDefault calls Get but returns dflt if the section or key was not found.
//
// Section and key strings are case-insensitive and ignore leading and trailing whitespace.
func (inif File) GetDefault(section, key, dflt string) string {
	return inif[Key(section)].GetDefault(key, dflt)
}
