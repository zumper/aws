package aws

type Element struct {
	Name      string `xml:"name,attr"`
	Type      string `xml:"type,attr"`
	MinOccurs string `xml:"minOccurs,attr"`
	MaxOccurs string `xml:"maxOccurs,attr"`
}

type ComplexType struct {
	Name    string    `xml:"name,attr"`
	Element []Element `xml:"sequence>element"`
	Group   GroupRef  `xml:"sequence>group"`          // Ignored for now
	Choice  []Choice  `xml:"sequence>choice>element"` // Ignored for now
}
