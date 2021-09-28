package main

import (
	"flag"
	"log"
	"log/syslog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lindsaybb/goPon"
)

var (
	helpFlag    = flag.Bool("h", false, "Show this help")
	verFlag     = flag.String("v", "1.5", "Version iter for reference")
	hostFlag    = flag.String("t", "10.5.100.10", "Hostname or IP Address of OLT")
	sleepFlag   = flag.Int("s", 10, "Sleep Interval for Loop")
	onceFlag    = flag.Bool("o", false, "Run once then exit")
	dregFlag    = flag.Bool("dr", false, "Don't deregister devices (double negative)")
	intfFlag    = flag.String("i", "", "Filter by interface (ex. 0/7)")
	localLog    = flag.Bool("ll", false, "Log Locally")
	syslogFlag  = flag.String("sl", "10.5.100.5:514", "Syslog Server (IP:Port)")
	pathFlag    = flag.String("p", "", "Path to create Logfile in")
	serviceFlag = flag.String("sp", "102_DATA_Acc", "Service Profile to Register Devices with")
	madeChange bool
)

func main() {
	flag.Parse()
	if *helpFlag {
		flag.PrintDefaults()
		os.Exit(0)
	}
	if *localLog {
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
	} else {
		slog, err := syslog.Dial("udp", *syslogFlag, syslog.LOG_INFO, "WainsWorld")
		if err != nil {
			log.Printf("FAILED to connect to syslog at %s\n", *syslogFlag)
		} else {
			log.SetOutput(slog)
		}
		defer slog.Close()
	}
	// create OLT object and verify connection
	olt := goPon.NewLumiaOlt(*hostFlag)
	if !olt.HostIsReachable() {
		log.Fatalln("Unreachable OLT")
	}
	// load currently registered device list into memory
	err := olt.UpdateOnuRegistry()
	if err != nil {
		log.Fatalln(err)
	}
	var obll *goPon.OnuBlacklistList
        var oil *goPon.OnuInfoList
        deregChan := make(chan *goPon.OnuInfo, 100)
        if !*dregFlag {
                go deregOnuCounter(olt, deregChan)
        }

	for {
		s := time.Now()
		// BlackList Registration Section
		// Will register any ONU with the Service Profile flag
		obll, err = olt.GetOnuBlacklist()
		if err != nil {
			log.Printf("Get Blacklist Error: %v\n", err)
		} else if len(obll.Entry) > 0 {
			for _, obl := range obll.Entry {
				// if Interface filter is supplied and doesn't match intf prefix
				if *intfFlag != "" && *intfFlag != obl.IfName {
					continue
				}
				newIntf, err := autoRegisterOnu(olt, obl)
				if err == nil {
					log.Printf("REGISTERED %s from Blacklist to %s\n", obl.SerialNumber, newIntf)
				} else {
					log.Printf("Error registering %s:%s\n", obl.SerialNumber, err)
				}
				madeChange = true
                                err = olt.UpdateOnuRegistry()
                                if err != nil {
                                        log.Printf("Error updating OLT Registry: %v\n", err)
                                }
			}
		}
		// Operationally Down ONU Remover
		// Will remove all ONU that shows as OperState != 1
		oil, err = olt.GetOnuInfoList()
		if err != nil {
			log.Printf("Error getting the ONU Info List: %v\n", err)
		} else if !*dregFlag {
			for _, e := range oil.Entry {
				if *intfFlag != "" && !strings.HasPrefix(e.IfName, *intfFlag) {
					continue
				}
				// checks if OperState == 1
				if !e.IsUp() {
					// send it to the dereg function for review
					deregChan <- e
				}
			}
		}
		if madeChange {
			olt.TabwriteRegistry()
			madeChange = false
		} else {
			// Sleep Flag says each loop will take a certain number of seconds minimum
			// if the loop completes without any activities, it should sleep
			// if the loop has activties, it should skip the sleep section
			si := time.Since(s)
			if si < (time.Duration(*sleepFlag) * time.Second) {
				time.Sleep((time.Duration(*sleepFlag) * time.Second) - si)
			}
		}
		if *onceFlag {
			return
		}
	}
}

func deregOnuCounter(olt *goPon.LumiaOlt, onu chan *goPon.OnuInfo) {
	dereg := make(map[string]int)
	for {
		// receive a new OnuInfo object
		entry := <-onu
		dereg[entry.SerialNumber]++
		if dereg[entry.SerialNumber] >= 10 {
			err := olt.DeauthOnuBySn(entry.SerialNumber)
			if err != nil {
				log.Printf("FAILED to deauth %s:%s\n%v\n", entry.SerialNumber, entry.IfName, err)
			} else {
				log.Printf("REMOVED Operationally Down ONU %s:%s\n", entry.SerialNumber, entry.IfName)
			}
			madeChange = true
			err = olt.UpdateOnuRegistry()
			if err != nil {
				log.Printf("Error updating OLT Registry: %v\n", err)
			}
			// reset list entry
			dereg[entry.SerialNumber] = 0
		} else {
			log.Printf("%s flap buffer: %d\n", entry.SerialNumber, dereg[entry.SerialNumber])
		}
	}

}

func autoRegisterOnu(olt *goPon.LumiaOlt, obl *goPon.OnuBlacklist) (string, error) {
	// create ONU Register Object from Blacklist info
	oreg := &goPon.OnuRegister{
		SerialNumber: obl.SerialNumber,
	}
	// update object with Interface to be registered on (next free)
	oreg = olt.NextAvailableOnuInterfaceUpdateRegister(obl.IfName, oreg)
	// use the oreg object to create an ONU Config
	onuCfg := goPon.NewOnuConfig(oreg.SerialNumber, oreg.Interface)
	// Override Authorize does not check Auth List or Vendor Prefix
	err := olt.AuthorizeOnuOverride(onuCfg)
	if err != nil {
		return oreg.Interface, err
	}
	// create ONU Service Profile according to supplied flag
	onuSvc := goPon.NewOnuProfile(onuCfg.IfName, *serviceFlag)
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

