
package vcloudoslicenses 

import (
    "encoding/xml"
    "net/url"
)

// QueryResultsRecords stuff below.

// <AdminVAppRecord vdcName="HMRC-MCTP-QA-IL2-STD-F-2 (IL2-PROD-STANDARD)" vdc="https://api.vcd.portal.skyscapecloud.com/api/vdc/71535bab-2731-42d8-84e2-2582573b3681" storageKB="41943040" status="POWERED_ON" ownerName="436.62.027659" org="https://api.vcd.portal.skyscapecloud.com/api/org/60a0dcac-f5c6-4e4d-bdbb-3a06586011b3" numberOfVMs="1" name="qa-app-1" memoryAllocationMB="16384" isVdcEnabled="true" isInMaintenanceMode="false" isExpired="false" isEnabled="true" isDeployed="true" isBusy="false" creationDate="2013-12-16T14:58:45.553Z" cpuAllocationMhz="4" href="https://api.vcd.portal.skyscapecloud.com/api/vApp/vapp-0022a90d-2535-4a12-8115-7c6e84e1d988" honorBootOrder="true" pvdcHighestSupportedHardwareVersion="8" cpuAllocationInMhz="8000" taskStatus="success" lowestHardwareVersionInVApp="8" task="https://api.vcd.portal.skyscapecloud.com/api/task/06a7386e-aca9-4a34-a56c-f06d9bf17926" numberOfCpus="4" taskStatusName="vappDeploy"/>

type AdminVAppRecord struct {
    XMLName string `xml:"AdminVAppRecord"`

    VDCName             string `xml:"vdcName,attr"`
    VDC                 string `xml:"vdc,attr"`
    StorageKB           string `xml:"storageKB,attr"`
    Status              string `xml:"status,attr"`
    OwnerName           string `xml:"ownerName,attr"`
    Org                 string `xml:"org,attr"`
    NumberOfVMs         int    `xml:"numberOfVMs,attr"`
    Name                string `xml:"name,attr"`
    MemoryAllocationMB  string `xml:"memoryAllocationMB,attr"`
    IsExpired           string `xml:"isExpired,attr"`
    IsEnabled           string `xml:"isEnabled,attr"`
    IsDeployed          string `xml:"isDeployed,attr"`
    CreationDate        string `xml:"creationDate,attr"`
    Href                string `xml:"href,attr"`
    NumberOfCpus        string `xml:"numberOfCpus,attr"`
}

// type VAppQueryResults struct {
//     XMLName string `xml:"AdminVAppRecord"`
// }

type VAppQueryResultsRecords struct {
    XMLName string `xml:"QueryResultRecords"`

    // TotalVapps int `xml:"total,attr"`

    Records []*AdminVAppRecord `xml:"AdminVAppRecord"`
}

// VDC VApp stuff below.

type VmOS struct {
    XMLName     string `xml:"OperatingSystemSection"`

    Id          string `xml:"id,attr"`
    Type        string `xml:"type,attr"`
    Href        string `xml:"href,attr"`
    OSType      string `xml:"osType,attr"`
}

type VAppVm struct {
    XMLName     string `xml:"Vm"`

    Deployed    string `xml:"deployed,attr"`
    Status      string `xml:"status,attr"`
    Name        string `xml:"name,attr"`
    Id          string `xml:"id,attr"`
    Type        string `xml:"type,attr"`
    Href        string `xml:"href,attr"`

    OperatingSystemSection *VmOS `xml:"OperatingSystemSection"`
}

type VAppChildren struct {
    XMLName     string `xml:"Children"`
    VM          []*VAppVm `xml:"Vm"`
}

type VDCVApp struct {
    XMLName     string `xml:"VApp"`

    Deployed    string `xml:"deployed,attr"`
    Status      string `xml:"status,attr"`
    Name        string `xml:"name,attr"`
    Id          string `xml:"id,attr"`
    Type        string `xml:"type,attr"`
    Href        string `xml:"href,attr"`

    VMs         VAppChildren `xml:"Children"`
}

func (a *VDCVApp) Get (session *VCloudSession, vdc *VdcResourceEntity) {
    r, _ := session.Get(vdc.Href)
    defer r.Body.Close()

    _ = xml.NewDecoder(r.Body).Decode(a)

    // if a.VMS != nil {
        for k1, v1 := range a.VMs.VM {
            u, _ := url.Parse(v1.Href)
            a.VMs.VM[k1].Href = u.Path 

            u, _ = url.Parse(v1.OperatingSystemSection.Href)
            a.VMs.VM[k1].OperatingSystemSection.Href = u.Path
        }
    // }
}

