
package vcloudoslicenses 

import (
    "strings"
    "sync"
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

type WorkerJob struct {
    Session         *vCloudSession
    Waiter          *sync.WaitGroup
    ResultsChannel  chan <- *ReportDocument
    Organisation    *OrganisationReference
}

func ReportWorker (job *WorkerJob) {
    vdcs := &VDCs{}
    vdcs.GetAll(job.Session, job.Organisation)

    for _, vdc := range vdcs.Records {
        if vdc.Type == "application/vnd.vmware.vcloud.vdc+xml" {
            vapps := &vApps{}
            vapps.GetAll(job.Session, vdc)   

            for _, vapp := range vapps.Records.Entities {
                if vapp.Type == "application/vnd.vmware.vcloud.vApp+xml" {
                    vms := &VMs{}
                    vms.GetAll(job.Session, vapp)

                    report := &ReportDocument{
                        Timestamp:      "NIL",
                        Year:           "NIL",
                        Month:          "NIL",
                        Day:            "NIL",
                        Organisation:   job.Organisation.Name,
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

                    job.ResultsChannel <- report
                }
            }
        }
    }

    job.Waiter.Done() 
}

func Report (session *vCloudSession) (report []*ReportDocument) {
    var reports []*ReportDocument
    var maxorgs int = 30

    waiter  := &sync.WaitGroup{}
    results := make(chan *ReportDocument)

    waiter.Add(maxorgs)

    orgs := &Organisations{}
    orgs.GetAll(session, "references", maxorgs)

    for _, org := range orgs.Records {
        job := &WorkerJob{
            Session:        session,
            Waiter:         waiter,
            ResultsChannel: results,
            Organisation:   org,
        }

        go ReportWorker(job)
    }

    go func() {
        waiter.Wait()
        close(results)
    }()

    for report := range results {
        reports = append(reports, report)
    }
  
    return reports 
}