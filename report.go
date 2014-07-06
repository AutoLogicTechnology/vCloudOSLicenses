
package vcloudoslicenses 

import (
    // "strings"
    "sync"
    // "time"
    // "strconv"

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
}

type WorkerJob struct {
    Waiter          *sync.WaitGroup
    ResultsChannel  chan <- *ReportDocument
    Organisation    *OrganisationReference
}

// func (v *VCloudSession) ReportWorker (job *WorkerJob) {
//     log.Print("Inside the worker...")

//     vdcs := &VDCs{}
//     vdcs.GetAll(v, job.Organisation)

//     for _, vdc := range vdcs.Records {
//         log.Print("Going over VDCs...")

//         if vdc.Type == "application/vnd.vmware.vcloud.vdc+xml" {
//             vapps := &vApps{}
//             vapps.GetAll(v, vdc)   

//             for _, vapp := range vapps.Records.Entities {
//                 log.Print("Going over vApps...")

//                 if vapp.Type == "application/vnd.vmware.vcloud.vApp+xml" {
//                     vms := &VMs{}
//                     vms.GetAll(v, vapp)

//                     now := time.Now()

//                     report := &ReportDocument{
//                         Timestamp:      now.String(),
//                         Year:           strconv.Itoa(now.Year()),
//                         Month:          now.Month().String(),
//                         Day:            strconv.Itoa(now.Day()),
//                         Organisation:   job.Organisation.Name,
//                         VDC:            vdc.Name,
//                         VApp:           vapp.Name,
//                         MSWindows:      0,
//                         RHEL:           0,
//                         CentOS:         0,
//                         Ubuntu:         0,
//                     }

//                     for _, vm := range vms.Records.Server {
//                         log.Print("And VMs...")

//                         if strings.Contains(vm.OSType.Name, "windows") {
//                             report.MSWindows++
//                         } else if strings.Contains(vm.OSType.Name, "rhel") {
//                             report.RHEL++
//                         } else if strings.Contains(vm.OSType.Name, "centos") {
//                             report.CentOS++
//                         } else if strings.Contains(vm.OSType.Name, "ubuntu") {
//                             report.Ubuntu++
//                         }
//                     }

//                     log.Printf("Report from worker: %+v", report)
//                     job.ResultsChannel <- report
//                 }
//             }
//         }
//     }

//     log.Print("Worker done...")
//     job.Waiter.Done() 
// }

func (v *VCloudSession) LicenseReport (max_organisations, max_pages int) (report []*ReportDocument) {
    var reports []*ReportDocument
    
    if max_organisations <= 0 {
        max_organisations = 10
    }

    if max_pages <= 0 {
        max_pages = 1
    }

    // waiter  := &sync.WaitGroup{}
    // results := make(chan *ReportDocument)

    // waiter.Add(max_organisations)

    orgs, _ := FindOrganisations(v, 10, 1)
    // orgs.GetAll(v, "references", max_organisations, max_pages)

    log.Printf("Orgs: %+v", orgs)

    // for _, org := range orgs.Records {
    //     job := &WorkerJob{
    //         Waiter:         waiter,
    //         ResultsChannel: results,
    //         Organisation:   org,
    //     }

    //     go v.ReportWorker(job)
    // }

    // go func() {
    //     waiter.Wait()
    //     close(results)
    // }()

    // for report := range results {
    //     reports = append(reports, report)
    // }
  
    return reports 
}