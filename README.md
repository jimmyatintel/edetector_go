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

### 1.0.0 (2023/10/02)
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
- Introduced Mock agent
  - Completed functions: handshake, detect, scan, collect, and explorer
  - Implemented a script to generate mock agents with random IDs, IPs, and MAC addresses

Fixed Bugs:
- Resolved the issue of generating incorrect new agent IDs (without "-")
- Enhanced task completion timing
- Let agent go offline when receiving GiveInfo (to create redis key in the beginning)

To-Do:
- Implement feature flags

### 1.0.1
Fixed bugs:
- Fixed empty columns error of white, black, and hack list