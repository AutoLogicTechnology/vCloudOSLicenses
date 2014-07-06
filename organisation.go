
package vcloudoslicenses 

import (
    "log"
    "fmt"

    "encoding/xml"
    "net/url"
)

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

    Records     []*OrgLink `xml:"Link"`
}
