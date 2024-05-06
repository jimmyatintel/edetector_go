# eDetector Server Documentation
Integrating with the agent for data reception, processing, and storage

This repository contains four microservices:

1. **Working Server**<br />
- Handle data from agents
- Send memory data to RabbitMQ

2. **DB Parser**<br />
- Parse Collect database
- Send Collect data to RabbitMQ

3. **Tree Builder**<br />
- Analyze relationships between files
- Send Explorer data to RabbitMQ

4. **Connector**<br />
- Send data to Elasticsearch
- 4 queues with different speed and tasks

## Getting Started
### Requirements
- Go Version: 1.20.8 linux/amd64
- 

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

## Directory Structure
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
  ├── dbRawData                     # Storage of finished DB files
  ├── ImageFile                     # Storage of finished Image files
  ├── *Staged                       # Files that have finished parsing
  ├── *Unstage                      # Files that are parsing
  ├── *Working                      # Files that are receiving data from agents
  ├── docs                          # Deprecated version documents
  ├── go.mod
  ├── go.sum
  ├── internal                      # Private library code
  ├── pkg                           # Public library code
  ├── README.md
  ├── static
  │   └── IP2LOCATION-LITE-DB5.BIN  # IP2LOCATION reference
  │   └── yaraRule                  # Files for YaraRule task
  │       ├── yaraNew.zip
  │       └── yara.zip
  └── test                          # Files for go test
  ```

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
  - Load different key image lists based on the OS
  - Use .tar.gz for compressing and decompressing file
  - Use "Ubuntu" as file system type instead of "Linux"

Fixed Bugs:
- Check TaskID when receiving finish signals to avoid inconsistent finish signals
- Show multiple sub roots in memory tree

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
- Memory Tree
- Key Image (New version)

Enhancements:
- Seamlessly replace data from the old to the new in the Collect and Explorer tasks
  - Add a finished signal to the RabbitMQ publisher at the end of the task
  - Upon receiving the finished signal, delete all old data by specifying the taskID
- Update the risk level function
  - Move the Hack List & VirusTotal check to the final stage to prevent its impact from being overwritten by other checks
- Store raw database data on the server
- Add new Collect tables
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
- Error creating new detectNetwork process
- Catch elastic response error
- Fix bugs of blacklist
- Fix online errors due to "NoKey" at the beginning.

### 1.0.3 (2023/12/18)
*Compatible Agent Version: Agent_1.0.5*

Enhancements:
- Introduce VirusTotal for network ip & process hash
  - Update riskScoure using VirusTotal
- Introduce FAT32
- Log the error when the scan crashes because of the agent
- Implement better terminate method for builder/parser
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
- Add Collect table: wireless
- Modify Collect table: process (add columns from the scan)
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
- Change the finish timing of Collect

### 1.0.0 (2023/10/05)
Finished tasks:
- Handshake
- Detect (memory & network)
- Scan (using zip files)
- Collect
- Explorer
- Image (temporary version)
- Agent update
- Terminate
- Graylog
- Log system
- White, black, and hack list (temporary version)

Enhancements:
- Introduce Mock agent
  - Complete functions: Handshake, Detect, Scan, Collect, and Explorer
  - Generate mock agents with random IDs, IPs, and MACs using go routine

Fixed Bugs:
- Fix the issue of generating incorrect new agent IDs (without "-")
- Enhance task completion timing
- Let agent go offline when receiving GiveInfo (to create redis key in the beginning)
- Fix empty columns error of white, black, and hack list
