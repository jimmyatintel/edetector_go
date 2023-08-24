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
| Agent → Server 	| GiveDetectNetwork     	| <details><summary>MemoryNetworkDetect</summary> ProcessId int json:"processId"<br>Address string json:"address"<br>Timestamp int json:"timestamp"<br>ProcessCreateTime int json:"processCreateTime"<br>ConnectionINorOUT bool json:"connectionInOrOut"<br>AgentPort int json:"agentPort"</details> |<details><summary>MemoryNetworkDetect</summary> 104984\|13.107.42.16:443\|1690922105\|1690080351\|0\|52365| |
| Server → Agent 	| DataRight             	|        	|          	|      	|
| Agent → Server 	| GiveDetectProcessFrag 	|   \(detect data fragment\)     	|          	|   Split into multiple fragments if it's too long   	|
| Server → Agent 	| DataRight             	|        	|          	|      	|
| Agent → Server 	| GiveDetectProcess     	|    <details><summary>Memory</summary> ProcessName string json:"processName"<br>ProcessCreateTime int json:"processCreateTime"<br>DynamicCommand string json:"dynamicCommand"<br>ProcessMD5 string json:"processMD5"<br>ProcessPath string json:"processPath"<br>ParentProcessId int json:"parentProcessId"<br>ParentProcessName string json:"parentProcessName"<br>ParentProcessPath string json:"parentProcessPath"<br>DigitalSign string json:"digitalSign"<br>ProcessId int json:"processId"<br>InjectActive string json:"injectActive"<br>ProcessBeInjected int json:"processBeInjected"<br>Boot string json:"boot"<br>Hide string json:"hide"<br>ImportOtherDLL string json:"importOtherDLL"<br>Hook string json:"hook"<br>ProcessConnectIP string json:"processConnectIP"<br>RiskLevel int json:"riskLevel"<br>Mode string json:"mode"<details>    	|          	|   Single or the last one   	|
| Server → Agent 	| DataRight             	|        	|          	|      	|
| Server → Agent 	| DataRight             	|        	|          	|      	|

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