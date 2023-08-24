# Working Server Documentation

## Functions
<details>
<summary>Handshake</summary>


| Direction      | TaskName            | Format                                                                                                                                                                     | Note |
| -------------- | ------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---- |
| Agent → Server | GiveInfo            | <details><summary>AgentInfo</summary> x64\|Windows 10 Home\|DESKTOP\-LD2C4NP\|SYSTEM\|3\.4\.2\.0,1988,1989\|20230815110126\|569a2191ae414802a5a72bc0b8e0bd1e\|0 </details> |      |
| Server → Agent | OpenCheckThread     | <details><summary>AgentID</summary> 6b75775ef8854658a595286f6f051399 </details>                                                                                            |      |
| Agent → Server | GiveDetectInfoFirst | <details><summary>process\|network</summary> 0\|0 </details>                                                                                                               |      |
| Server → Agent | UpdateDetectMode    | <details><summary>process\|network</summary> 0\|0 </details>                                                                                                               |      |
| Agent → Server | GiveDetectInfo      | <details><summary>process\|network</summary> 0\|0 </details>                                                                                                               |      |
| Server → Agent | CheckConnect        | \(Heartbeat\)                                                                                                                                                              |      |


</details>

<details>

<summary>ChangeDetectMode</summary>



</details>

<details>

<summary>Detect</summary>



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