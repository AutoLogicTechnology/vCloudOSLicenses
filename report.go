
package vcloudoslicenses 

import (
    "strings"
    "sync"
    "time"
    "strconv"
    "encoding/xml"
    "log"
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
    Unknown         uint
    TotalVMs        uint
}

type WorkerJob struct {
    WorkerID        int 
    Waiter          *sync.WaitGroup
    ResultsChannel  chan <- *ReportDocument
    RecycleChannel  chan <- *VAppQueryResultsRecords
    Organisation    *Organisation
    VApps           *VAppQueryResultsRecords
}

type VAppWorkerJob struct {
    Waiter          *sync.WaitGroup
    ResultsChannel  chan <- *ReportDocument
    VApp            *AdminVAppRecord
}

func (v *VCloudSession) ReportWorker (job *WorkerJob) {
    for _, link := range job.Organisation.Links {
        if link.Type == "application/vnd.vmware.vcloud.vdc+xml" {
            v.Counters.VDCs++

            vdc := &VDC{}
            vdc.Get(v, link)

            for _, entity := range vdc.ResourceEntities.ResourceEntity {
                if entity.Type == "application/vnd.vmware.vcloud.vApp+xml" {
                    v.Counters.VApps++

                    vapp := &VDCVApp{}
                    vapp.Get(v, entity)

                    now := time.Now()
                    report := &ReportDocument{
                        Timestamp:      now.String(),
                        Year:           strconv.Itoa(now.Year()),
                        Month:          now.Month().String(),
                        Day:            strconv.Itoa(now.Day()),
                        Organisation:   job.Organisation.Name,
                        VDC:            vdc.Name,
                        VApp:           vapp.Name,
                        MSWindows:      0,
                        RHEL:           0,
                        CentOS:         0,
                        Ubuntu:         0,
                        Unknown:        0,
                    }

                    for _, vm := range vapp.VMs.VM {
                        v.Counters.VMs++

                        if strings.Contains(vm.OperatingSystemSection.OSType, "windows") {
                            report.MSWindows++
                        } else if strings.Contains(vm.OperatingSystemSection.OSType, "rhel") {
                            report.RHEL++
                        } else if strings.Contains(vm.OperatingSystemSection.OSType, "centos") {
                            report.CentOS++
                        } else if strings.Contains(vm.OperatingSystemSection.OSType, "ubuntu") {
                            report.Ubuntu++
                        } else {
                            report.Unknown++
                        }
                    }

                    job.ResultsChannel <- report
                }
            }
        }
    }

    job.Waiter.Done() 
}

func (v *VCloudSession) VAppReportWorker (job *WorkerJob) {
    var recycled *VAppQueryResultsRecords = nil 

    for _, vapp := range job.VApps.Records {
        vdc := &VDCVApp{}

        r, err := v.Get(vapp.Href)

        if err != nil {
            if recycled == nil {
                log.Print("Creating new recycled bin...")
                recycled = &VAppQueryResultsRecords{}
            }

            log.Printf("Adding vApp to recycle bin: vapp not found? %s", vapp.Href)
            recycled.Records = append(recycled.Records, vapp)

            continue 
        }

        defer r.Body.Close()

        _ = xml.NewDecoder(r.Body).Decode(vdc)

        v.Counters.VApps++

        org := &Organisation{}
        err = org.Get(v, vapp.Org)
        if err != nil {
            if recycled == nil {
                log.Print("Creating new recycled bin...")
                recycled = &VAppQueryResultsRecords{}
            }

            log.Printf("Adding vApp to recycle bin: org not found? %s", vapp.Org)
            recycled.Records = append(recycled.Records, vapp)

            continue 
        }

        now := time.Now()
        report := &ReportDocument{
            Timestamp:      now.String(),
            Year:           strconv.Itoa(now.Year()),
            Month:          now.Month().String(),
            Day:            strconv.Itoa(now.Day()),
            Organisation:   org.Name,
            VDC:            vapp.VDCName,
            VApp:           vapp.Name,
            MSWindows:      0,
            RHEL:           0,
            CentOS:         0,
            Ubuntu:         0,
            Unknown:        0,
            TotalVMs:       vapp.NumberOfVMs,
        }

        for _, vm := range vdc.VMs.VM {
            v.Counters.VMs++

            if strings.Contains(vm.OperatingSystemSection.OSType, "windows") {
                report.MSWindows++
            } else if strings.Contains(vm.OperatingSystemSection.OSType, "rhel") {
                report.RHEL++
            } else if strings.Contains(vm.OperatingSystemSection.OSType, "centos") {
                report.CentOS++
            } else if strings.Contains(vm.OperatingSystemSection.OSType, "ubuntu") {
                report.Ubuntu++
            } else {
                report.Unknown++
            }
        }

        job.ResultsChannel <- report
    }

    if recycled != nil {
        log.Print("Passing recycle into recycle channel...")
        job.RecycleChannel <- recycled
    }

    job.Waiter.Done() 
}

func (v *VCloudSession) VAppReport (max_vapps, max_pages int) (reports []*ReportDocument) {
    var worker_id int = 1
    
    if max_vapps <= 0 {
        max_vapps = 10
    }

    if max_pages <= 0 {
        max_pages = 1
    }

    waiter  := &sync.WaitGroup{}
    
    results := make(chan *ReportDocument)
    recycled := make(chan *VAppQueryResultsRecords)

    vapps, _ := v.FindVApps(max_vapps, max_pages)

    waiter.Add(len(vapps))

    for _, vapp := range vapps {
        job := &WorkerJob{
            WorkerID:       worker_id,
            Waiter:         waiter,
            ResultsChannel: results,
            RecycleChannel: recycled,
            VApps:          vapp,
        }

        go v.VAppReportWorker(job)
        worker_id++
    }

    go func() {
        waiter.Wait()
        close(results)
        close(recycled)
    }()

    for report := range results {
        reports = append(reports, report)
    }

    log.Print("Going over recycle bin...")
    for orphen := range recycled {
        for _, vapp := range orphen.Records {
            log.Printf("Found recycled vApps: %+v", vapp.Name)
        }
    }
  
    return reports 
}

func (v *VCloudSession) LicenseReport (max_organisations, max_pages int) (report []*ReportDocument) {
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

    orgs, _ := v.FindOrganisations(max_organisations, max_pages)
    v.Counters.Orgs = len(orgs)

    for _, org := range orgs {
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