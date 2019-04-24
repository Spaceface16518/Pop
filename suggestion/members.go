package suggestion

// NewSuggestions is a
func NewSuggestions(initialData *map[string]int) Suggestions {
	s := Suggestions{}
	if initialData != nil {
		s.AddAll(initialData)
	}
	return s
}

// Add appends a name and count to a Suggestions struct. Use this method instead of appending manually as this method will be resistant to implementation change
func (suggestions *Suggestions) Add(name string, count int) {
	suggestions.Suggestions = append(suggestions.Suggestions, Suggestion{name, count})
}

// AddAll adds the contents of a map to an existing Suggestions struct. Internally uses Add.
func (suggestions *Suggestions) AddAll(suggestionsMap *map[string]int) {
	for k, v := range *suggestionsMap {
		suggestions.Add(k, v)
	}
}
