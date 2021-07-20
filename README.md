# wainsWorld

wainsWorld is an Iskratel PON utility program that uses the lindsaybb/gopon library for automatic registration of ONU with a data service profile, and includes deregistration of the ONU when it becomes operationally down. This program has passed through 500 Huawei devices in Lindsay's production test environment to fulfill a strand-mounted PON gateway order. Further updates are anticipated as we continue with 4000 more pieces, 2000 of which are ZTE. Dedicated to our production manager, Wayne, the steam engine.

| Flag | Description |
| ------ | ------ |
| -h | Show this help |
| -t | Hostname or IP Address of OLT (default "10.5.100.10" for our production OLT) |
| -s | Sleep interval for Loop (default 5, interpreted as seconds) |
| -o | Run program once then exit (default false) |
| -i | Filter by interface (example "0/7" will only reg/dereg devices on port 7) |
| -p | Path to create Logfile in (folder, filename is auto-generated) |
| -sp | Service Prifle to Register Devices with (default "102_DATA_Acc" for ONU UNI 1 Gigabit service) |

# Example Call
```sh
./wainsWorld -p ~/prod >> ~/prod/so37543.log &
tail -f ~/prod/2021-07-20T11-01-10.log
```

# Example Output
```sh
------------
GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanOnuCfgTable
------------
------------
GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanServicePortProfileTable
------------
------------
PATCH Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanOnuCfgTable/msanOnuCfgEntry=0%2F7%2F8
------------
------------
POST Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanServicePortProfileTable/msanServicePortProfileEntry=0%2F7%2F8
------------
------------
GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanOnuCfgTable
------------
------------
GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanServicePortProfileTable
------------
------------
GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanOnuInfoTable
------------
------------
DELETE Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanServicePortProfileTable/msanServicePortProfileEntry=0%2F1%2F9,101_CWMP
------------
------------
PATCH Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanOnuCfgTable/msanOnuCfgEntry=0%2F1%2F9
------------
------------
GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanOnuCfgTable
------------
------------
GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanServicePortProfileTable
------------
------------
DELETE Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanServicePortProfileTable/msanServicePortProfileEntry=0%2F1%2F13,101_CWMP
------------
------------
PATCH Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanOnuCfgTable/msanOnuCfgEntry=0%2F1%2F13
------------
------------
GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanOnuCfgTable
------------
------------
GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanServicePortProfileTable
------------
------------
DELETE Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanServicePortProfileTable/msanServicePortProfileEntry=0%2F1%2F14,101_CWMP
------------
------------
PATCH Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanOnuCfgTable/msanOnuCfgEntry=0%2F1%2F14
------------
------------
GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanOnuCfgTable
------------
------------
GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanServicePortProfileTable
------------

```

# Example Logfile
```sh
2021/07/20 11:09:28 REGISTERED ZTEGC663A981 from Blacklist to 0/7/8
2021/07/20 11:09:56 REMOVED Operationally Down ONU ZTEGC661A3C5:0/7/6
2021/07/20 11:10:56 REGISTERED ZTEGC667A24B from Blacklist to 0/7/6
2021/07/20 11:11:46 REMOVED Operationally Down ONU ZTEGC6675BA8:0/7/7
2021/07/20 11:12:57 REGISTERED ZTEGC6650574 from Blacklist to 0/7/7
2021/07/20 11:13:25 REMOVED Operationally Down ONU ZTEGC66E09D7:0/7/5
2021/07/20 11:13:48 REMOVED Operationally Down ONU ZTEGC6650574:0/7/7
2021/07/20 11:14:27 REGISTERED ZTEGC6650574 from Blacklist to 0/7/5
2021/07/20 11:15:17 REMOVED Operationally Down ONU ZTEGC66565C3:0/7/4
2021/07/20 11:15:44 REGISTERED ZTEGC66EA6BD from Blacklist to 0/7/4
2021/07/20 11:16:13 REMOVED Operationally Down ONU ZTEGC66EA6BD:0/7/4
2021/07/20 11:16:51 REGISTERED ZTEGC66EA6BD from Blacklist to 0/7/4
2021/07/20 11:17:20 REMOVED Operationally Down ONU ZTEGC66EA6BD:0/7/4
2021/07/20 11:17:48 REGISTERED ZTEGC662C26E from Blacklist to 0/7/4
2021/07/20 11:18:16 REMOVED Operationally Down ONU ZTEGC66453AD:0/7/3
2021/07/20 11:18:39 REMOVED Operationally Down ONU ZTEGC662C26E:0/7/4
2021/07/20 11:19:07 REGISTERED ZTEGC667D0D8 from Blacklist to 0/7/3
2021/07/20 11:19:30 REGISTERED ZTEGC66EA6BD from Blacklist to 0/7/4
2021/07/20 11:19:58 REMOVED Operationally Down ONU ZTEGC666A456:0/7/2
2021/07/20 11:20:26 REGISTERED ZTEGC662C26E from Blacklist to 0/7/2
2021/07/20 11:21:00 REGISTERED ZTEGC6656BE4 from Blacklist to 0/7/7
2021/07/20 11:21:28 REMOVED Operationally Down ONU ZTEGC66BA643:0/7/1
2021/07/20 11:22:28 REGISTERED ZTEGC669207D from Blacklist to 0/7/1
2021/07/20 11:22:57 REMOVED Operationally Down ONU ZTEGC663A981:0/7/8
2021/07/20 11:23:52 REMOVED Operationally Down ONU ZTEGC667A24B:0/7/6
2021/07/20 11:24:19 REGISTERED ZTEGC6631A8B from Blacklist to 0/7/6
2021/07/20 11:24:48 REMOVED Operationally Down ONU ZTEGC6631A8B:0/7/6
2021/07/20 11:25:16 REGISTERED ZTEGC666C53B from Blacklist to 0/7/6
2021/07/20 11:25:44 REMOVED Operationally Down ONU ZTEGC6650574:0/7/5
2021/07/20 11:26:07 REMOVED Operationally Down ONU ZTEGC666C53B:0/7/6
2021/07/20 11:26:35 REGISTERED ZTEGC6691A01 from Blacklist to 0/7/5
2021/07/20 11:26:58 REGISTERED ZTEGC6631A8B from Blacklist to 0/7/6
2021/07/20 11:27:26 REMOVED Operationally Down ONU ZTEGC66EA6BD:0/7/4
2021/07/20 11:27:54 REGISTERED ZTEGC666C53B from Blacklist to 0/7/4
2021/07/20 11:28:23 REMOVED Operationally Down ONU ZTEGC662C26E:0/7/2
2021/07/20 11:28:45 REMOVED Operationally Down ONU ZTEGC666C53B:0/7/4
2021/07/20 11:29:13 REGISTERED ZTEGC66A4C75 from Blacklist to 0/7/2
2021/07/20 11:29:41 REMOVED Operationally Down ONU ZTEGC667D0D8:0/7/3
2021/07/20 11:30:09 REGISTERED ZTEGC666C53B from Blacklist to 0/7/3

```
