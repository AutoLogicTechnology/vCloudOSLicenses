
package main 

import (
    "log"
    "strings"
)

type ReportDocument struct {
    Timestamp       string 
    Year            string 
    Month           string 
    Day             string 
    Organisation    string
    VDC             string 
    VApp            string
    MSWindows       uint 
    RHEL            uint 
    CentOS          uint 
    Ubuntu          uint 
}

func ReportWorker (id int, session *vCloudSession, jobs <- chan *OrganisationReference, results chan <- *ReportDocument) {

    log.Printf("Inside Worker: %d", id)

    for org := range jobs {
        vdcs := &VDCs{}
        vdcs.GetAll(session, org)
        
        if len(vdcs.Records) <= 0 {
            continue 
        }

        for _, vdc := range vdcs.Records {
            if vdc.Type == "application/vnd.vmware.vcloud.vdc+xml" {
                vapps := &vApps{}
                vapps.GetAll(session, vdc)   

                for _, vapp := range vapps.Records.Entities {
                    if vapp.Type == "application/vnd.vmware.vcloud.vApp+xml" {
                        // log.Printf("vApp: %v", vapp.Name)

                        vms := &VMs{}
                        vms.GetAll(session, vapp)

                        report := &ReportDocument{
                            Timestamp:      "NIL",
                            Year:           "NIL",
                            Month:          "NIL",
                            Day:            "NIL",
                            Organisation:   org.Name,
                            VDC:            vdc.Name,
                            VApp:           vapp.Name,
                            MSWindows:      0,
                            RHEL:           0,
                            CentOS:         0,
                            Ubuntu:         0,
                        }

                        for _, vm := range vms.Records.Server {
                            if strings.Contains(vm.OSType.Name, "windows") {
                                report.MSWindows++
                            } else if strings.Contains(vm.OSType.Name, "rhel") {
                                report.RHEL++
                            } else if strings.Contains(vm.OSType.Name, "centos") {
                                report.CentOS++
                            } else if strings.Contains(vm.OSType.Name, "ubuntu") {
                                report.Ubuntu++
                            }
                        }

                        log.Printf("vApp: %s: Windows = %d, RHEL = %d, CentOS = %d, Ubuntu = %d", vapp.Name, report.MSWindows, report.RHEL, report.CentOS, report.Ubuntu)
                        results <- report 
                    }
                }
            }
        }
    }
}

func Report (session *vCloudSession) (report []*ReportDocument) {
    log.SetFlags(log.Lmicroseconds)

    var reports []*ReportDocument
    var maxorgs int = 5

    jobs := make(chan *OrganisationReference, 500)
    results := make(chan *ReportDocument, 5000)

    for i := 1; i <= maxorgs; i++ {
        go ReportWorker(i, session, jobs, results)
    }

    orgs := &Organisations{}
    orgs.GetAll(session, "references", maxorgs)

    log.Printf("Orgs Found: %d", len(orgs.Records))

    for _, org := range orgs.Records {
        log.Print("Adding Organisation to job queue")

        jobs <- org 
    }

    close(jobs)

    for report := range results {
        reports = append(reports, report)
    }

    close(results)

    return reports 
}