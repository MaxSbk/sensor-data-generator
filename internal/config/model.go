package config

type Config struct {
	Mqtt struct {
		Url      string `yaml:"url" env:"MQTT_URL" env-default:"tcp://localhost:1883"`
		ClientId string `yaml:"client-id"`
		Topic    string `yaml:"topic" env:"MQTT_TOPIC" env-default:"sensor-data"`
		Qos      byte   `yaml:"qos"`
	} `yaml:"mqtt"`

	Sensors []Sensor `yaml:"sensors"`
}

type Sensor struct {
	Id        string    `yaml:"id"`
	Type      string    `yaml:"type"`
	MachineId string    `yaml:"machine-id"`
	PartId    string    `yaml:"part-id"`
	ToolId    string    `yaml:"tool-id"`
	Unit      string    `yaml:"unit"`
	Generator Generator `yaml:"generator"`
}

type Generator struct {
	Range            Range      `yaml:"values"`
	Interval         uint       `yaml:"interval"`
	ExtraBelowValues ExtraValue `yaml:"extra_below_values"`
	ExtraAboveValues ExtraValue `yaml:"extra_above_values"`
}

type Range struct {
	Min float64 `yaml:"min"`
	Max float64 `yaml:"max"`
}

type ExtraValue struct {
	Freq                uint    `yaml:"freq"`
	PercentageDeviation float64 `yaml:"percentage_deviation"`
	Duration            uint    `yaml:"duration"`
}
