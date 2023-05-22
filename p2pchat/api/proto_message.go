package api

type Message []Segment

type Segment struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type TextSegment struct {
	Text string `json:"text"`
}

func (msg Message) ExtractPlaintext() string {
	var text string
	for _, segment := range msg {
		if segment.Type == "Text" {
			data := segment.Data.(TextSegment)
			text += data.Text
		}
	}
	return text
}

func PlaintextToMessage(text string) Message {
	return Message{Segment{Type: "text", Data: TextSegment{Text: text}}}
}
