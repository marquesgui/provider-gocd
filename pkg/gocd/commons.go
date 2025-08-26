package gocd

// ConfigProperty represents a key/value property.
type ConfigProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (c ConfigProperty) Equal(o ConfigProperty) bool {
	return c.Key == o.Key && c.Value == o.Value
}

// HALLinks represents standard GoCD _links
type HALLinks struct {
	Self *HALLink `json:"self,omitempty"`
	Doc  *HALLink `json:"doc,omitempty"`
	Find *HALLink `json:"find,omitempty"`
}

type HALLink struct {
	Href string `json:"href"`
}
