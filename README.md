# vCloud OS Licenses

A small library for calculating the OS types, and their numbers, within a vCloud Director setup. Loops over the organisations, finding the OrgVDCs within each one, and in turn, finding each vApp within the VDC. Finally, the process ends with a loop over the VMs in the vApp, tallying up the OS types.

The minimum possible amount of work is done with the XML returned by the vCloud API. Please don't expect to find a full vCloud XML to Go struct mapping here.