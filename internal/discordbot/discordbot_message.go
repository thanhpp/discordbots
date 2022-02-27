package discordbot

import "strings"

type Message struct {
	Topic   string
	Content []*MsgField
}

type MsgField struct {
	Key string
	Val string
}

func (m *Message) Stringtify() string {
	var strBuilder = new(strings.Builder)

	strBuilder.WriteString("\n")
	strBuilder.WriteString(">>> **")
	strBuilder.WriteString(strings.ToUpper(m.Topic))
	strBuilder.WriteString("** \n")

	for _, field := range m.Content {
		strBuilder.WriteString("`")
		strBuilder.WriteString(field.Key)
		strBuilder.WriteString("`")
		strBuilder.WriteString(": ")
		strBuilder.WriteString(field.Val)
		strBuilder.WriteString("\n")
	}

	return strBuilder.String()
}

func (m *Message) AddContent(key, val string) {
	m.Content = append(m.Content, &MsgField{
		Key: key,
		Val: val,
	})
}
