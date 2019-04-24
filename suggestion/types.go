package suggestion

// Suggestion is a struct for html templating, consisting of the name associated with a suggestion and it's upvote count. The count may be negative.
type Suggestion struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// Suggestions is a simple struct that wraps a list of suggestions for feeding into and html template
type Suggestions struct {
	Suggestions []Suggestion `json:"suggestions"`
}

// TODO: keep in certain order
// TODO: Use just this instead of a map[string]int
