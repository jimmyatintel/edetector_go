# eDetector Server Documentation
Integrating with agents for receiving, processing, and storing data.

## Introduction

### Task List

- 記憶體偵測：Scan
- 記憶體掃描：DetectProcess & DetectNetwork
- 記憶體樹：MemoryTree
- 痕跡取證：Collection
- 檔案總表：Explorer
- 關鍵映像檔：KeyImage
- YaraRule
- 任務終止：Termination
- Agent 更新：UpdateAgent
- Agent 移除：RemoveAgent

### Flow

1. Establish connections with services (MariaDB, Redis, RabbitMQ, and ElasticSearch), the API server, and agents
2. Receive tasks from the API server
3. Trigger agents with the tasks
4. Process raw data from agents
5. Make data available to users through the following methods:
   - Directly insert or update data in ElasticSearch
   - Send data to RabbitMQ and then to ElasticSearch
   - Store data for user download (e.g. KeyImage)

### Notion Documents

- [The flow, formats, and details of all tasks](https://www.notion.so/edetector/Working-Server-Functions-Doc-e4ea043d66b44ad484ee0b172281b7f2?pvs=4)
- [Meaning of agent task status and progress](https://www.notion.so/edetector/Agent-Task-Status-Progress-Doc-421a1a2501b841ec93402d4d0a52d655?pvs=4)
- [Version Dependency with other services](https://www.notion.so/edetector/Version-Doc-27b63115e38f44d3afaf5fceb68c2eec?pvs=4)

### Microservices

Please enable the following four microservices and ensure they run continuously to ensure the execution of tasks.

1. **Working Server**<br />
    - Receive all tasks from the API server 
    - Receive and process raw data from agents
    - Send data of Scan, DetectProcess, DetectNetwork, and MemoryTree to RabbitMQ
      - Calculate Risk score and level
    - Store Collection and Explorer files in "dbUstaged" and "fileUstaged" (They will be used by other microservices)
    - Store KeyImage files
    - Update YaraRule information to the Explorer

2. **DB Parser**<br />
    - Parse Collection files in "dbUstaged"
      - Use sqlite3 analyzer
      - Transfer all time formats to Unix timestamp
      - Insert data corresponding to their table names
    - Send Collection data to RabbitMQ

3. **Tree Builder**<br />
    - Build Explorer Tree using files from "fileUnstage"
      - Record relationships between files and find the root directory
      - Traverse the tree to generate the full path for each file
    - Send Explorer data to RabbitMQ

4. **Connector**<br />
    - Bulk Insert data to Elasticsearch
      - Use four queues with different speeds and tasks

### Directory Structure
  ```
  .
  ├── cmd                           # Entry point of the services
  │   ├── builder
  │   │   ├── builder.go            # Entry point of the builder
  │   │   └── log_builder.log       # Log file for the builder
  │   ├── connector
  │   │   ├── connector.go          # Entry point of the connector
  │   │   └── log_conn.log          # Log file for the connector
  │   ├── mockagent
  │   │   ├── agent.go
  │   │   └── log_agent.log
  │   ├── parser
  │   │   ├── log_parser.log        # Entry point of the parser
  │   │   └── parser.go             # Log file for the parser
  │   └── server
  │       ├── log_gin.log           # Log file for gin
  │       ├── log_server.log        # Log file for server
  │       └── working_server.go     # Entry point of the server
  ├── commit_example.md
  ├── config                        # Config files
  │   ├── app.env
  │   ├── config.go
  │   └── container.yml
  ├── go.mod
  ├── go.sum
  ├── internal                      # Private library code
  ├── pkg                           # Public library code
  ├── README.md
  ├── static
  |   ├── dbRawData                     # Storage of finished DB files
  |   ├── ImageFile                     # Storage of finished Image files
  |   ├── *Staged                       # Files that have finished parsing
  |   ├── *Unstage                      # Files that are parsing
  |   ├── *Working                      # Files that are receiving data from agents
  │   ├── IP2LOCATION-LITE-DB5.BIN  # IP2LOCATION reference
  │   └── yaraRule                  # Files for YaraRule task
  │       ├── yaraNew.zip
  │       └── yara.zip
  └── test                          # Files for go test
  ```

## Getting Started
### Requirements
- Go Version: 1.20.8 linux/amd64

### Installation
```bash
cd edetector_go
go mod download
```

### Run the Services
```bash
go run cmd/server/working_server.go
go run cmd/connector/connector.go
go run cmd/builder/builder.go
go run cmd/parser/parser.go
```

### Monitor Logs
You can view service logs from the console, `./cmd/<service>/*.log`, and `/var/log/syslog`

## Version Documentation
### 1.1.1
*Compatible Agent Version: Agent_1.0.9*

Enhancements:
- Introduce the new elasticsearch version: 7.17
  - Remove "_type" : "_doc" when inserting data
- Store WARN level log to mariaDB
- Introduce the new version of yara rule task
  - Remove the part of GivePath in yara rule task
  - Use New format for matched rules
- Check the OS of agents
- Introduce the Linux(Ubuntu) agent
  - Load different KeyImage lists based on the OS
  - Use .tar.gz for compressing and decompressing file
  - Use "Ubuntu" as file system type instead of "Linux"
- Use the WebSocket for updating task status and agent connection status
- Move "working", "unstage", and "staged" directories to the static directory
- Use "COLLECT_READY" and "IMAGE_ERROR" functions to detect crash

Fixed Bugs:
- Check TaskID when receiving finish signals to avoid inconsistent finish signals
- Show multiple sub roots in MemoryTree

### 1.1.0 (2024/04/29)
*Compatible Agent Version: Agent_1.0.9*

Enhancements:
- Introduce the new Elasticsearch structure by merging indices from 59 to 3
  - Adjust the code according to the new nested mapping format.
- Implement a better method for opening files in the tree builder
- Implement a retry mechanism when Elasticsearch crashes

### 1.0.5 (2024/04/11)
*Compatible Agent Version: Agent_1.0.8*

New Tasks:
- MemoryTree
- KeyImage (New version)

Enhancements:
- Seamlessly replace data from the old to the new in the Collection and Explorer tasks
  - Add a finished signal to the RabbitMQ publisher at the end of the task
  - Upon receiving the finished signal, delete all old data by specifying the taskID
- Update the risk level function
  - Move the Hack List & VirusTotal check to the final stage to prevent its impact from being overwritten by other checks
- Store raw database data on the server
- Add new Collection tables
  - Handle more various time formats
- Control the number of agents
  - Add AGENT_LIMIT and TOTAL_AGENT_LIMIT env variables
  - Offline all agents at the beginning and end
- Use a better split format
  - Utilize |@| instead of | in GiveInfo, Scan, DetectProcess
- Implement version dependency between agent and server
  - Add MIN_AGENT_VERSION env variable

Fixed Bugs:
- Delete inserted data upon task failure

### 1.0.4 (2024/02/19)
*Compatible Agent Version: Agent_1.0.5*

Enhancements:
- Introduce YaraRule
- New task type (terminating, different error types)

Fixed Bugs:
- Error creating new DetectNetwork process
- Catch elastic response error
- Fix bugs of blacklist
- Fix online errors due to "NoKey" at the beginning.

### 1.0.3 (2023/12/18)
*Compatible Agent Version: Agent_1.0.5*

Enhancements:
- Introduce VirusTotal for network ip & process hash
  - Update riskScoure using VirusTotal
- Introduce FAT32
- Log the error when the Scan crashes because of the agent
- Implement better Termination method for builder/parser
  - add new status "terminating"
- Implement go unit test
- Add RejectAgent TaskType

Fixed Bugs:
- Fix the issue of Error storing log to database: Data too long for column 'content'
- Fix bugs of parser format
- Fix the deadlock issues in rbconnector service
- Connect MariaDB in rbconnector service
- Change the log level of "Invalid line" from 'Warn' to 'Error'

Working in progress:
- Implement feature flags


### 1.0.2 (2023/11/27)
*Compatible Agent Version: Agent_1.0.4*

Enhancements:
- Add Collection table: wireless
- Modify Collection table: process (add columns from Scan)
- Introduce RemoveAgent
- Introduce .tar.gz for linux agents
- Introduce ip2location
- Adjust prefetchCount to rabbitMQ

Fixed Bugs:
- Remove unnecessary panic()
- Fix read buffer error
- Close unused go routines

### 1.0.1 (2023/10/12)
Enhancements:
- Check the data format using the number of columns
- Mock agent
  - Added image task
  - Sleep between each iteration

Fixed Bugs:
- GiveInfo steps: mySQL -> redis -> request
- Change the finish timing of Collection

### 1.0.0 (2023/10/05)
Finished tasks:
- Handshake
- DetectProcess & DetectNetwork
- Scan (using zip files)
- Collection
- Explorer
- Image (temporary version)
- UpdateAgent
- Termination
- Graylog
- Log system
- White, black, and hack list (temporary version)

Enhancements:
- Introduce Mock agent
  - Complete functions: Handshake, DetectProcess, DetectNetwork, Scan, Collection, and Explorer
  - Generate mock agents with random IDs, IPs, and MACs using go routine

Fixed Bugs:
- Fix the issue of generating incorrect new agent IDs (without "-")
- Enhance task completion timing
- Let agent go offline when receiving GiveInfo (to create redis key in the beginning)
- Fix empty columns error of white, black, and hack list
