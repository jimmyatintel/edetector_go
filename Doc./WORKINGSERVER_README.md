# Working Server Documentation

## Functions
<details>
<summary>Handshake</summary>

<style>
    .task {
        width: 80%;
    }
    .task th {
        word-wrap: break-word;
    }
    .task tr:nth-child(even) {
        background-color: gray;
    }
</style>

<div class="task">

| Direction      | TaskName            | Format           | Msg e\.g\. | Note |
|----------------|---------------------|------------------|------------|------|
| Agent → Server | GiveInfo            | \(Agent info\)   |            |      |
| Server → Agent | OpenCheckThread     | \(Agent id\)     |            |      |
| Agent → Server | GiveDetectInfoFirst | process\|network | 0\|0       |      |
| Server → Agent | UpdateDetectMode    | process\|network | 0\|0       |      |
| Agent → Server | GiveDetectInfo      | process\|network | 0\|0       |      |
| Server → Agent | CheckConnect        | \(Heartbeat\)    |            |      |


</div>

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