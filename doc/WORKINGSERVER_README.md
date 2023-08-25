# Working Server Documentation

## Functions
<details>
<summary>Handshake</summary>

| Direction      | TaskName            | Format           | Msg e.g.                                                                                                                                                                                | Note |
| -------------- | ------------------- | ---------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---- |
| Agent → Server | GiveInfo            | AgentInfo        | <nobr><details><summary>AgentInfo</summary> x64\|Windows 10 Home\|DESKTOP\-LD2C4NP\|SYSTEM\|3\.4\.2\.0,1988,1989\|20230815110126\|569a2191ae414802a5a72bc0b8e0bd1e\|0 </details></nobr> |      |
| Server → Agent | OpenCheckThread     | AgentID          | <nobr><details><summary>AgentID</summary> 6b75775ef8854658a595286f6f051399 </details></nobr>                                                                                            |      |
| Agent → Server | GiveDetectInfoFirst | process\|network | 0\|0                                                                                                                                                                                    |      |
| Server → Agent | UpdateDetectMode    | process\|network | 0\|0                                                                                                                                                                                    |      |
| Agent → Server | GiveDetectInfo      | process\|network | 0\|0                                                                                                                                                                                    |      |
| Server → Agent | CheckConnect        | \(Heartbeat\)    | 0\|0                                                                                                                                                                                    |      |

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

| Direction      | TaskName              | Format                                      | Msg e.g.                                                                                                                          | Note                                                        |
| -------------- | --------------------- | ------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------- |
| Agent → Server | GiveDetectNetwork     | <nobr>\{struct MemoryNetworkDetect\}</nobr> | <nobr><details><summary>MemoryNetworkDetect</summary> 104984\|13.107.42.16:443\|1690922105\|1690080351\|0\|52365</details></nobr> |                                                             |
| Server → Agent | DataRight             |                                             |                                                                                                                                   |                                                             |
| Agent → Server | GiveDetectProcessFrag | <nobr>\(detect data fragment\)</nobr>       |                                                                                                                                   | <nobr>Split into multiple fragments if it's too long</nobr> |
| Server → Agent | DataRight             |                                             |                                                                                                                                   |                                                             |
| Agent → Server | GiveDetectProcess     | <nobr>\{struct Memory\}</nobr>              | <nobr><details><summary>Memory</summary></details></nobr>                                                                         | <nobr>Single or the last one</nobr>                         |
| Server → Agent | DataRight             |                                             |                                                                                                                                   |                                                             |
| Server → Agent | DataRight             |                                             |                                                                                                                                   |                                                             |

<blockquote>
<details><summary><em>MemoryNetworkDetect</em></summary>
    ProcessId <code>int</code><br>
    Address <code>string</code><br>
    Timestamp <code>int</code><br>
    ProcessCreateTime <code>int</code><br>
    ConnectionINorOUT <code>bool</code><br>
    AgentPort <code>int 
</details>
</blockquote>

<blockquote>
<details><summary><em>Memory</em></summary>
    ProcessName <code>string</code><br>
    ProcessCreateTime <code>int</code><br>
    DynamicCommand <code>string</code><br>
    ProcessMD5 <code>string</code><br>
    ProcessPath <code>string </code><br>
    ParentProcessId <code>int </code><br>
    ParentProcessName <code>string</code><br>
    ParentProcessPath <code>string</code><br>
    DigitalSign <code>string</code><br>
    ProcessId <code>int</code><br>
    InjectActive <code>string</code><br>
    ProcessBeInjected <code>int</code><br>
    Boot <code>string</code><br>
    Hide <code>string</code><br>
    ImportOtherDLL <code>string</code><br>
    Hook <code>string</code><br>
    ProcessConnectIP <code>string</code><br>
    RiskLevel <code>int</code><br>
    Mode <code>string</code>
</details>
</blockquote>

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