package utils

type YAMLRawMessage struct {
	unmarshal func(interface{}) error
}

func (msg *YAMLRawMessage) UnmarshalYAML(unmarshal func(interface{}) error) error {
	msg.unmarshal = unmarshal
	return nil
}

func (msg *YAMLRawMessage) Unmarshal(v interface{}) error {
	return msg.unmarshal(v)
}
