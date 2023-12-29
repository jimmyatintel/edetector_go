# eDetector Server Documentation
Integrating with the agent for data reception, processing, and storage

This repository contains four microservices : 

1. **Working Server**<br />
- Handle data from agents
- Send detect & scan data (memory) to RabbitMQ

2. **DB Parser**<br />
- Parse collect database
- Send collect data to RabbitMQ

3. **Tree Builder**<br />
- Analyze relationships between files
- Send explorer data to RabbitMQ

4. **Connector**<br />
- Send data to Elasticsearch
- 3 queues with different speed

## Version

### 1.0.0 (2023/10/05)
Finished tasks:
- Handshake
- Detect (process & memory)
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
  - Complete functions: handshake, detect, scan, collect, and explorer
  - Generate mock agents with random IDs, IPs, and MACs using go routine

Fixed Bugs:
- Fix the issue of generating incorrect new agent IDs (without "-")
- Enhance task completion timing
- Let agent go offline when receiving GiveInfo (to create redis key in the beginning)
- Fix empty columns error of white, black, and hack list

### 1.0.1 (2023/10/12)
Enhancements:
- Check the data format using the number of columns
- Mock agent
  - Added image task
  - Sleep between each iteration

Fixed Bugs:
- GiveInfo steps: mySQL -> redis -> request
- Change the finish timing of collect

### 1.0.2 (2023/11/27)
Compatible Agent Version: Agent_1.0.4

Enhancements:
- Add collect table: wireless
- Modify collect table: process (add columns from the scan)
- Introduce RemoveAgent
- Introduce .tar.gz for linux agents
- Introduce ip2location
- Adjust prefetchCount to rabbitMQ

Fixed Bugs:
- Remove unnecessary panic()
- Fix read buffer error
- Close unused go routines

### 1.0.3 (2023/12/18)
Compatible Agent Version: Agent_1.0.5

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

### 1.0.4
Compatible Agent Version: Agent_1.0.5

Fixed Bugs:
- Error creating new detectNetwork process
- Catch elastic response error

## Directory Structure
```
.
├── agentFile                     # Agent executions files for updating
│   └── agent.exe
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
├── dbStaged                      # DB files that have finished parsing
├── dbUnstage                     # DB files that are parsing
├── dbWorking                     # DB files that are receiving data from agents
├── docs
├── fileStaged                    # File txt files that have finished building
├── fileUnstage                   # File txt files that are bulilding
├── fileWorking                   # File txt files that are receiving data from agents
├── go.mod
├── go.sum
├── imageUnstage                  # Image files that have finished received
├── imageWorking                  # Image files that are receiving data from agents
├── internal                      # Private library code
├── mockFiles                     # Mock agent files
├── pkg                           # Public library code
├── README.md
├── scanUnstage                   # Scan txt files that have finished received
├── scanWorking                   # Scan txt files that are receiving data from agents
├── static
│   └── IP2LOCATION-LITE-DB5.BIN  # IP2LOCATION reference
└── test                          # Files for go test
```
