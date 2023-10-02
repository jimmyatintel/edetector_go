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

### 1.0.0 (2023/9/23)
Finished tasks:
- Handshake
- Detect (process & memory)
- Scan (using zip files)
- Collect
- Explorer
- Terminate
- Graylog
- Log system

ToDo:
- Image task
- Agent update task
- White, black, and hack list
- Feature flag

### 1.0.1 (2023/10/02)
Enhancements:
- Introduced Mock agent
  - Completed functions: handshake, detect, scan, collect, and explorer
  - Implemented a script to generate mock agents with random IDs, IPs, and MAC addresses
- Introduced Image task (temporary version)
- Introduced Agent update task
- Implemented white, black, and hack lists

Fixed Bugs:
- Resolved the issue of generating incorrect new agent IDs (without "-")
- Enhanced task completion timing

To-Do:
- Implement feature flags