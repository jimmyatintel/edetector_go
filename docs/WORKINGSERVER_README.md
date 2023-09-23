# Working Server Documentation

## Version

### 2.0.7 ()
Enhancements:
- Add "Terminate" mission
- Add "Image" mission
- Add "ReadyScan" function
- Enhance readability for "fromclient.go"
- Identify the task type from agent connection
- Elevate certain error-level messages to the panic level
- Zip scan data

### 2.0.6 (2023/8/31)
Enhancements:
- Add "Terminate" mission
- Merge internal/taskservice/taskdb.go into pkg/mariadb/query/taskservice.go
- Move touser.go to pkg

### 2.0.5 (2023/8/30)
Enhancements:
- Enhance readability for the init function
- Update progress without notifying frontend
  - Add environment variable: UPDATE_INTERVAL
- Improve Naming & file structure: sendToElastic -> ToRabbitMQ
- Display the version at the start
- Simplify "SendTCPtoClient"
- Move file & fflag from internal to pkg
- Merge risklevel & memorystruct into work package

Fixed bugs:
- Invoke "RequestToUser" through a regular call instead of using a go routine

### 2.0.4 (2023/8/29)
Fixed bugs:
- Fix incorrect entry point of "updateDriveProgress"

### 2.0.3 (2023/8/28)
Enhancements:
- Use environment variables for the first progress part

### 2.0.2 (2023/8/28)
Fixed bugs:
- Fix error parsing memory_network_scan (incorrect entry point of "scanNetworkElastic")

### 2.0.1 (2023/8/28)
Enhancements:
- Redirect the Gin logs to Graylog
- Add scan network to elasticsearch

Fixed bugs:
- Map failed subtasks to their corresponding tasks in MariaDB
- Set the entry points of collect progress
- Remove the "GiveCollectInfo" function
- Fix incorrect key parameter in the "memory_network_detect"
- Skip FAT drives

### 2.0.0 (2023/8/24)
Description:
- Fully functionality for the **new agent**
- Redirect the log to the graylog

### 1.0.3 (2023/8/28)
Enhancements:
- Redirect the Gin logs to Graylog

### 1.0.2 (2023/8/24)
Enhancements:
- Redirect the log to the graylog

### 1.0.1 (2023/8/17)
Description:
- Fully functionality for the **old agent**
