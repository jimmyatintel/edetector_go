# eDetector Server Documentation
This repository contains four microservices : 

- **Working Server**<br />

- **DB Parser**<br />

- **Tree Builder**<br />

- **Connector**<br />

## Version

### 2.0.2 (2023/8/28)
Enhancements:
- Use environment variables for the first progress part
Fixed bugs:
- Fix bugs in the "memory_network_scan" (error parsing)

### 2.0.1 (2023/8/28)
Enhancements:
- Redirect the Gin logs to Graylog
- Add ack and count limit to the connector to prevent excessive consumption
- Add scan network to the elasticsearch

Fixed bugs:
- Mapp failed subtasks to their corresponding tasks in MariaDB
- Set the entry points of progress collection
- Remove the "GiveCollectInfo" function
- Fix bugs in the "memory_network_detect" (incorrect key parameter)
- Skip FAT drives

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