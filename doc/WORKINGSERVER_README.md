# Working Server Documentation

## Functions
<details>
<summary>Handshake</summary>

| Direction      | TaskName            | Format           | Msg e.g.                                                                                                                                                                   | Note |
| -------------- | ------------------- | ---------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---- |
| Agent → Server | GiveInfo            | AgentInfo        | <details><summary>AgentInfo</summary> x64\|Windows 10 Home\|DESKTOP\-LD2C4NP\|SYSTEM\|3\.4\.2\.0,1988,1989\|20230815110126\|569a2191ae414802a5a72bc0b8e0bd1e\|0 </details> |      |
| Server → Agent | OpenCheckThread     | AgentID          | <details><summary>AgentID</summary> 6b75775ef8854658a595286f6f051399 </details>                                                                                            |      |
| Agent → Server | GiveDetectInfoFirst | process\|network | 0\|0                                                                                                                                                                       |      |
| Server → Agent | UpdateDetectMode    | process\|network | 0\|0                                                                                                                                                                       |      |
| Agent → Server | GiveDetectInfo      | process\|network | 0\|0                                                                                                                                                                       |      |
| Server → Agent | CheckConnect        | \(Heartbeat\)    | 0\|0                                                                                                                                                                       |      |

</details>

<details>

<summary>ChangeDetectMode</summary>

| Direction      	| TaskName         	| Format           	| Msg e.g. 	| Note 	|
|----------------	|------------------	|------------------	|----------	|------	|
| User → Server  	| ChangeDetectMode 	| process\|network 	| 0\|0     	|      	|
| Server → Agent 	| UpdateDetectMode 	| process\|network 	| 0\|0     	|      	|
| Agent → Server 	| GiveDetectInfo   	| process\|network 	| 0\|0     	|      	|

</details>

<details>

<summary>Detect</summary>

|      Direction 	| TaskName              	| Format 	| Msg e.g. 	| Note 	|
|---------------	|-----------------------	|--------	|----------	|------	|
| Agent → Server 	| GiveDetectNetwork     	|  \{MemoryNetworkDetect\} struct |<details><summary>MemoryNetworkDetect</summary> 104984\|13.107.42.16:443\|1690922105\|1690080351\|0\|52365| |
| Server → Agent 	| DataRight             	|        	|          	|      	|
| Agent → Server 	| GiveDetectProcessFrag 	|   \(detect data fragment\)     	|          	|   Split into multiple fragments if it's too long   	|
| Server → Agent 	| DataRight             	|        	|          	|      	|
| Agent → Server 	| GiveDetectProcess     	| \{Memory struct\}	|<details><summary>Memory</summary> 	|   Single or the last one   	|
| Server → Agent 	| DataRight             	|        	|          	|      	|
| Server → Agent 	| DataRight             	|        	|          	|      	|

<details><summary>MemoryNetworkDetect</summary> ProcessId int<br>Address string<br>Timestamp int<br>ProcessCreateTime int<br>ConnectionINorOUT bool <br>AgentPort int </details>
<details><summary>Memory</summary> ProcessName string<br>ProcessCreateTime int<br>DynamicCommand string<br>ProcessMD5 string<br>ProcessPath string <br>ParentProcessId int <br>ParentProcessName string<br>ParentProcessPath string<br>DigitalSign string<br>ProcessId int<br>InjectActive string<br>ProcessBeInjected int<br>Boot string<br>Hide string<br>ImportOtherDLL string<br>Hook string <br>ProcessConnectIP string<br>RiskLevel int<br>Mode string

</details>

<details>

<summary>Scan</summary>



</details>

<details>

<summary>Explorer</summary>



</details>

<details>

<summary>Collect</summary>



</details>

## Version

### 2.0.0 (2023/8/24)
Description:
- Fully functionality for the **new agent**
- Redirecting the log to the graylog

### 1.0.2 (2023/8/24)
Enhancements:
- Redirecting the log to the graylog

### 1.0.1 (2023/8/17)
Description:
- Fully functionality for the **old agent**