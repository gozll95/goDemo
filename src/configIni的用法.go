package main

import (
	"flag"
	"log"
	"runtime"

	"strings"

	"github.com/larspensjo/config"
)

var (
	configFile = flag.String("configfile", "config.ini", "General configuration file")
	hostFile   = flag.String("hostFile", "host", "target configuration file")
	varsFile   = flag.String("varsFile", "port_vars.yml", "target configuration file")
	pureFile   = flag.String("pureFile", "pure_vars.yml", "target configuration file")
)

func main() {
	//解析参数
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	cfg, err := config.ReadDefault(*configFile)
	if err != nil {
		log.Fatalf("Fail to find", *configFile, err)
	}

	//start role
	Role = ParseSection(cfg, "role")
	//fmt.Printf("Role is %v\n", Role)

	//start service
	ServiceHost = ParseSection(cfg, "service")
	//	fmt.Printf("ServiceHost is %v\n", ServiceHost)

	//start ip
	HostIp = ParseSection(cfg, "ip")
	//fmt.Printf("HostIp is %v\n", HostIp)

	//start PortTcp
	PortTcp = ParseSection(cfg, "port_tcp")
	//fmt.Printf("PortTcp is %v\n", PortTcp)

	//start PortHttp
	PortHttp = ParseSection(cfg, "port_http")
	//fmt.Printf("PortHttp is %v\n", PortHttp)

	makeServiceData(cfg)

	cw := config.NewDefault()
	cw1 := config.NewDefault()
	cw2 := config.NewDefault()

	makeNewIni(cw)
	cw.WriteFile(*hostFile, 0644, "Test file for test-case")

	makeVarsIni(cw1)
	cw1.WriteFile(*varsFile, 0644, "Test file for test-case")

	makePureVarsIni(cw2)
	cw2.WriteFile(*pureFile, 0644, "Test file for test-case")

	//makeVarsIni2(cw1)

}

func makeServiceData(cfg *config.Config) {
	for service, host := range ServiceHost {

		ServiceMetric = make(map[string][]string)
		//判断该service是否有metric
		if JudgeMetric(cfg, service) {
			ServiceMetric = ParseSection(cfg, service+"_metric")
		}
		//fmt.Printf("%s metric %v\n", service, ServiceMetric)
		/*
			ServiceHost is map[kodo_memcache:[role_meta] kodo_mysql:[zhu3] kodo_acc:[zhu1 zhu2] kodo_nginx:[web] kodo_mongodb:[role_meta]]
		*/
		KodoService[service] = &Service{}
		machine := Machine{}
		for _, singleHost := range host {
			//,oop singleHost in zhu1,zhu2,role_meta
			if strings.HasPrefix(singleHost, "role_") {
				roleHostLists := Role[strings.Trim(singleHost, "role_")]
				//loop roleHostLists in zhu1,zhu2,zhu3
				for _, sHost := range roleHostLists {
					//如果zhu1有metric的话
					if metric, ok := ServiceMetric[sHost]; ok {
						machine = Machine{
							Host:   sHost,
							Ip:     HostIp[sHost][0],
							Metric: ParseMetric(metric),
						}
						KodoService[service].HostLists = append(KodoService[service].HostLists, machine)

					} else {
						machine = Machine{
							Host: sHost,
							Ip:   HostIp[sHost][0],
						}
						KodoService[service].HostLists = append(KodoService[service].HostLists, machine)
					}
				}
			} else {
				machine = Machine{
					Host: singleHost,
					Ip:   HostIp[singleHost][0],
				}
				KodoService[service].HostLists = append(KodoService[service].HostLists, machine)
			}

		}
		if port, ok := PortTcp[service]; ok {
			KodoService[service].Port = port
			KodoService[service].Tag = "tcp"
		} else if port, ok := PortHttp[service]; ok {
			KodoService[service].Port = port
			KodoService[service].Tag = "http"
		}
		//fmt.Printf("%s is %+v\n", service, KodoService[service])
	}

}

func makeNewIni(cw *config.Config) {
	for service, v := range KodoService {
		for _, machine := range v.HostLists {
			var metric string
			for key, value := range machine.Metric {
				metric = metric + key + "=" + value + " "
			}
			cw.AddOption(service, machine.Host, metric)
		}
	}

}

func makeVarsIni(cw *config.Config) {
	//var ipList []string
	serviceIpList := make(map[string][]string)
	for service, v := range KodoService {
		var pre string
		switch v.Tag {
		case "http":
			pre = "http://"
		default:
			pre = ""
		}
		if len(v.Port) >= 2 {
			for _, port := range v.Port {
				var ipList []string
				for _, machine := range v.HostLists {
					ip := pre + machine.Ip + ":" + port
					ipList = append(ipList, ip)
				}
				serviceIpList[service+"_"+port] = ipList
				//ipList = []string{""}
			}
		} else if len(v.Port) == 1 {
			var ipList []string
			for _, machine := range v.HostLists {
				ip := pre + machine.Ip + ":" + v.Port[0]
				ipList = append(ipList, ip)
			}
			serviceIpList[service] = ipList
			//ipList = []string{""}
		} else {
			var ipList []string
			for _, machine := range v.HostLists {
				ip := pre + machine.Ip
				ipList = append(ipList, ip)
			}
			serviceIpList[service] = ipList
			//ipList = []string{""}
		}
	}

	for service, ipList := range serviceIpList {
		var ipString string
		ipString = "["
		//fmt.Println(service, len(ipList))
		for _, ip := range ipList[0 : len(ipList)-1] {
			{
				tmp := `"` + ip + `"` + ","
				ipString = ipString + tmp
			}
		}
		ipString = ipString + `"` + ipList[len(ipList)-1] + `"` + "]"
		cw.AddOption("", service, ipString)
	}
}

//	fmt.Println(serviceIpList)

func makePureVarsIni(cw *config.Config) {
	serviceIpList := make(map[string][]string)
	for service, v := range KodoService {
		var ipList []string
		for _, machine := range v.HostLists {
			ipList = append(ipList, machine.Ip)
		}
		serviceIpList["pure_"+service] = ipList
	}
	for service, ipList := range serviceIpList {
		var ipString string
		ipString = "["
		//fmt.Println(service, len(ipList))
		for _, ip := range ipList[0 : len(ipList)-1] {
			{
				tmp := `"` + ip + `"` + ","
				ipString = ipString + tmp
			}
		}
		ipString = ipString + `"` + ipList[len(ipList)-1] + `"` + "]"
		cw.AddOption("", service, ipString)
	}
}
