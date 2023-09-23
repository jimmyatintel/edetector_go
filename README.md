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
- Image
- Agent update
- White, black, and hack list