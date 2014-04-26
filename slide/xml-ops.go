package aws

type Operation struct {
	Name   string     `xml:"name,attr"`
	Input  InputElem  `xml:"input"`
	Output OutputElem `xml:"output"`
}

type InputElem struct {
	Message string `xml:"message,attr"`
}

type OutputElem struct {
	Message string `xml:"message,attr"`
}
