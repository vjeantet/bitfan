package api

// Doc represents a processor documentation
//
// A Doc is ....
//
// A Doc can have.....
//
// swagger:model processorDoc
type processorDoc struct {
	Kind     string `json:"kind"`
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

// Pipeline represents a pipeline
//
// A Pipeline is ....
//
// A Pipeline can have.....
//
// swagger:model Pipeline
type Pipeline struct {
	// the id for this pipeline
	ID int `json:"id"`
	// the Label
	// min length: 3
	Label string `json:"label"`
	// the location
	ConfigLocation string `json:"config_location"`
	// the location's host
	ConfigHostLocation string `json:"config_host_location"`

	Content string `json:"config_content"`
}

// Error represents a error
//
// A Error is ....
//
// A Error can have.....
//
// swagger:model Error
type Error struct {
	Message string `json:"error"`
}
