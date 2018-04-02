package models

// Doc represents a processor documentation
//
// A Doc is ....
//
// A Doc can have.....
//
// swagger:model processorDoc
type processorDoc struct {
	Kind     string `json:"kind"`
	Behavior string `json:"behavior"`
	Name     string `json:"name"`
	Doc      string `json:"doc"`
	DocShort string `json:"doc_short"`
	Options  *struct {
		Doc     string `json:"doc"`
		Options []*struct {
			Name         string      `json:"name"`
			Alias        string      `json:"alias"`
			Doc          string      `json:"doc"`
			Required     bool        `json:"requiered"`
			Type         string      `json:"type"`
			DefaultValue interface{} `json:"default_value"`
			//LogstashExample
			ExampleLS string `json:"example"`
		} `json:"options"`
	} `json:"options"`
	Ports []*struct {
		Default bool   `json:"default"`
		Name    string `json:"name"`
		Number  int    `json:"number"`
		Doc     string `json:"doc"`
	} `json:"ports"`
}
