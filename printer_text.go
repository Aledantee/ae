package ae

import "strings"

func (p *Printer) printsText(err error, depth int) string {
	var sb strings.Builder

	// TODO: implement text printing
	sb.WriteString(Message(err))

	return sb.String()
}
