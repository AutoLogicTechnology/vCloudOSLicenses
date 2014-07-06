
package vcloudoslicenses 

type OrgLink struct {
    XMLName     string `xml:"Link"`

    Type        string `xml:"type,attr"` 
    Name        string `xml:"name,attr"`
    Href        string `xml:"href,attr"`
}

type Organisation struct {
    XMLName     string `xml:"Org"`

    Name        string `xml:"name,attr"`
    Id          string `xml:"id,attr"`
    Type        string `xml:"type,attr"`
    Href        string `xml:"href,attr"`

    Links       []*OrgLink `xml:"Link"`
}
