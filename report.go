
package vcloudoslicenses 

import (
    "strings"
    "sync"
    "time"
    "strconv"
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
    Waiter          *sync.WaitGroup
    ResultsChannel  chan <- *ReportDocument
    Organisation    *OrganisationReference
}

func (v *VCloudSession) ReportWorker (job *WorkerJob) {
    vdcs := &VDCs{}
    vdcs.GetAll(v, job.Organisation)

    for _, vdc := range vdcs.Records {
        if vdc.Type == "application/vnd.vmware.vcloud.vdc+xml" {
            vapps := &vApps{}
            vapps.GetAll(v, vdc)   

            for _, vapp := range vapps.Records.Entities {
                if vapp.Type == "application/vnd.vmware.vcloud.vApp+xml" {
                    vms := &VMs{}
                    vms.GetAll(v, vapp)

                    now := time.Now()

                    report := &ReportDocument{
                        Timestamp:      now.String(),
                        Year:           strconv.Itoa(now.Year()),
                        Month:          strconv.Itoa(now.Month()),
                        Day:            strconv.Itoa(now.Day()),
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

func (v *VCloudSession) Report (max_organisations, max_pages int) (report []*ReportDocument) {
    var reports []*ReportDocument
    
    if max_organisations <= 0 {
        max_organisations = 10
    }

    if max_pages <= 0 {
        max_pages = 1
    }

    waiter  := &sync.WaitGroup{}
    results := make(chan *ReportDocument)

    waiter.Add(max_organisations)

    orgs := &Organisations{}
    orgs.GetAll(v, "references", max_organisations, max_pages)

    for _, org := range orgs.Records {
        job := &WorkerJob{
            Waiter:         waiter,
            ResultsChannel: results,
            Organisation:   org,
        }

        go v.ReportWorker(job)
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