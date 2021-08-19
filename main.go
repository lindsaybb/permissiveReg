package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lindsaybb/gopon"
)

var (
	helpFlag    = flag.Bool("h", false, "Show this help")
	verFlag     = flag.String("v", "1.4", "Version iter for reference")
	hostFlag    = flag.String("t", "10.5.100.10", "Hostname or IP Address of OLT")
	sleepFlag   = flag.Int("s", 5, "Sleep Interval for Loop")
	graceFlag   = flag.Int("g", 5, "Grace timer between reg & dereg intervals")
	onceFlag    = flag.Bool("o", false, "Run once then exit")
	dregFlag    = flag.Bool("dr", false, "Don't deregister devices (double negative)")
	intfFlag    = flag.String("i", "", "Filter by interface (ex. 0/7)")
	pathFlag    = flag.String("p", "", "Path to create Logfile in")
	serviceFlag = flag.String("sp", "102_DATA_Acc", "Service Profile to Register Devices with")
)

func main() {
	flag.Parse()
	if *helpFlag {
		flag.PrintDefaults()
		return
	}
	// generate date/time stamp
	dtStamp := newDateTimeStamp()
	dtStamp += ".log"
	lfPath := filepath.Join(*pathFlag, dtStamp)
	lfAbs, err := filepath.Abs(lfPath)
	if err != nil {
		log.Fatalln(err)
	}
	// create new log file
	logFile, err := os.OpenFile(lfAbs, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("FAILED to create logfile at %s\n", lfPath)
	} else {
		log.SetOutput(logFile)
	}
	defer logFile.Close()
	// create OLT object and verify connection
	olt := gopon.NewLumiaOlt(*hostFlag)
	if !olt.HostIsReachable() {
		log.Fatalln("Unreachable OLT")
	}
	// load currently registered device list into memory
	err = olt.UpdateOnuRegistry()
	if err != nil {
		log.Fatalln(err)
	}

	for {
		// BlackList Registration Section
		// Will register any ONU with the Service Profile flag
		obll, err := olt.GetOnuBlacklist()
		if err != nil {
			log.Printf("Get Blacklist Error: %v\n", err)
		} else if len(obll.Entry) > 0 {
			for _, obl := range obll.Entry {
				newIntf, err := autoRegisterOnu(olt, obl)
				if err == nil {
					log.Printf("REGISTERED %s from Blacklist to %s\n", obl.SerialNumber, newIntf)
					err = olt.UpdateOnuRegistry()
					if err != nil {
						log.Printf("Error updating OLT Registry: %v\n", err)
					}
				} else {
					log.Printf("Error registering %s:%s\n", obl.SerialNumber, err)
				}
			}
		}
		time.Sleep(time.Duration(*graceFlag) * time.Second)
		// Operationally Down ONU Remover
		// Will remove all ONU that shows as OperState != 1
		oil, err := olt.GetOnuInfoList()
		if err != nil {
			log.Printf("Error getting the ONU Info List: %v\n", err)
		} else if !*dregFlag {
			for _, e := range oil.Entry {
				//fmt.Printf("%s:%v\n", e.IfName, e.IsUp())
				// first check if interface is in filter range
				// allows for bad filters to restrict program operation
				if *intfFlag != "" && !strings.HasPrefix(e.IfName, *intfFlag) {
					continue
				}
				// checks if OperState == 1
				if !e.IsUp() {
					// MIGHT WANT TO check again or wait to avoid flapping
					// this function also removes all service profiles
					err = olt.DeauthOnuBySn(e.SerialNumber)
					if err != nil {
						log.Printf("FAILED to deauth %s:%s\n%v\n", e.SerialNumber, e.IfName, err)
					} else {
						log.Printf("REMOVED Operationally Down ONU %s:%s\n", e.SerialNumber, e.IfName)
						err = olt.UpdateOnuRegistry()
						if err != nil {
							log.Printf("Error updating OLT Registry: %v\n", err)
						}
					}

				}
			}
		}
		olt.TabwriteRegistry()
		if *onceFlag {
			return
		}
		time.Sleep(time.Duration(*sleepFlag) * time.Second)
	}
}

func autoRegisterOnu(olt *gopon.LumiaOlt, obl *gopon.OnuBlacklist) (string, error) {
	// create ONU Register Object from Blacklist info
	oreg := &gopon.OnuRegister{
		SerialNumber: obl.SerialNumber,
	}
	// update object with Interface to be registered on (next free)
	oreg = olt.NextAvailableOnuInterfaceUpdateRegister(obl.IfName, oreg)
	// use the oreg object to create an ONU Config
	onuCfg := gopon.NewOnuConfig(oreg.SerialNumber, oreg.Interface)
	// Override Authorize does not check Auth List or Vendor Prefix
	err := olt.AuthorizeOnuOverride(onuCfg)
	if err != nil {
		return oreg.Interface, err
	}
	// create ONU Service Profile according to supplied flag
	onuSvc := gopon.NewOnuProfile(onuCfg.IfName, *serviceFlag)
	// if service does not exist, ONU will still register but not pass data
	err = olt.PostOnuProfile(onuSvc)
	return oreg.Interface, err
}

func newDateTimeStamp() string {
	t := time.Now()
	// formatting the time so it can be used in path without spaces or colons
	tf := strings.Replace(t.Format(time.RFC3339), ":", "-", -1) // replace : in time format with -
	return tf[:len(tf)-6]                                       // remove trailing nanoseconds
}

