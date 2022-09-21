package internals

type Plan struct {
	FormatVersion      string           `json:"format_version"`
	TerraformVersion   string           `json:"terraform_version"`
	ResourceChanges    []ResourceChange `json:"resource_changes"`
	PriorState         State            `json:"prior_state"`
	PlannedValue       Values           `json:"planned_values"`
	ProposedUnknown    Values           `json:"proposed_unknown"`
	RelevantAttributes []Attribute      `json:"relevant_Attributes"`
	Configuration      Configuration    `json:"configuration"`
}

type State struct {
	Values           Values `json:"values"`
	TerraformVersion string `json:"terraform_version"`
}

type Change struct {
	Actions      []string    `json:"actions"`
	Before       interface{} `json:"before"`
	After        interface{} `json:"after"`
	AfterUnknown interface{} `json:"after_unknown"`
	ReplacePaths [][]string  `json:"replace_paths,omitempty"`
}

type ResourceChange struct {
	Address         string `json:"address"`
	PreviousAddress string `json:"previous_address"`
	ModuleAddress   string `json:"module_address"`
	Type            string `json:"type"`
	Name            string `json:"name"`
	Provider        string `json:"provider_name"`
	Changes         Change `json:"change"`
	Reason          string `json:"action_reason"`
}

type Values struct {
	RootModule RootModule             `json:"root_module"`
	Outputs    map[string]OutputValue `json:"outputs"`
}

type RootModule struct {
	Resources    []Resource    `json:"resources"`
	ChildModules []ChildModule `json:"child_modules"`
}

type ChildModule struct {
	Address      string        `json:"address"`
	Resources    []Resource    `json:"resources"`
	ChildModules []ChildModule `json:"child_modules"`
}

type Resource struct {
	Address         string                 `json:"address"`
	Mode            string                 `json:"mode"`
	Type            string                 `json:"type"`
	Name            string                 `json:"name"`
	Index           interface{}            `json:"index"`
	ProviderName    string                 `json:"provider_name"`
	SchemaVersion   int                    `json:"schema_version"`
	Values          map[string]interface{} `json:"values"`
	SensitiveValues interface{}            `json:"sensitive_values,omitempty"`
}

type OutputValue struct {
	Value     interface{} `json:"value"`
	Type      string      `json:"type"`
	Sensitive bool        `json:"sensitive"`
}
type Attribute struct {
	Resource  string `json:"resource"`
	Attribute string `json:"attribute"`
}

type ProviderConfigs map[string]ProviderConfig

type Configuration struct {
	ProviderConfig ProviderConfigs `json:"provider_config"`
	Root           ModuleConfig    `json:"root_module"`
}

type ProviderConfig struct {
	Name        string      `json:"name"`
	FullName    string      `json:"full_name"`
	Alias       string      `json:"alias,omitempty"`
	Address     string      `json:"module_address"`
	Expressions interface{} `json:"expressions"`
}

type Expression struct {
	Constant   interface{} `json:"constant_value,omitempty"`
	References []string    `json:"references"`
}

type BlockExpression map[string]Expression

type ResourceConfig struct {
	Address           string        `json:"address"`
	Mode              string        `json:"mode"`
	Type              string        `json:"type"`
	Name              string        `json:"name"`
	ProviderConfigKey string        `json:"provider_config_key"`
	Provisioners      []Provisioner `json:"provisioners"`
	Expressions       interface{}   `json:"expressions"`
	SchemaVersion     int           `json:"schema_version"`
	CountExpression   Expression    `json:"count_expression"`
	ForEachExpression Expression    `json:"for_each_expression"`
}

type Provisioner struct {
	Type        string          `json:"type"`
	Expressions BlockExpression `json:"expressions"`
}
type OutputConfig struct {
	Expression Expression `json:"expression"`
	Sensitive  bool       `json:"sensitive"`
}

type ModuleConfig struct {
	Outputs   map[string]OutputConfig `json:"outputs"`
	Resources []ResourceConfig        `json:"resources"`
	Modules   ModuleCalls             `json:"module_calls"`
}

type ModuleCallChild struct {
	ResolvedSource    string          `json:"resolved_source"`
	Expressions       BlockExpression `json:"expressions"`
	CountExpression   Expression      `json:"count_expression"`
	ForEachExpression Expression      `json:"for_each_expression"`
	Module            ModuleConfig    `json:"module"`
}

type ModuleCalls map[string]ModuleCallChild
