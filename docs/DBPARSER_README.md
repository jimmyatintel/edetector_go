# DB Parser Documentation

## Version

### 2.0.3 ()
Enhancements:
- Add extension name to the log to differentiate between the builder and parser
- Elevate certain error-level messages to the panic level

### 2.0.3 (2023/8/31)
Enhancements:
- Merge internal/taskservice/taskdb.go into pkg/mariadb/query/taskservice.go

### 2.0.2 (2023/8/30)
Enhancements:
- Enhance readability for the init function
- Improve Naming & file structure: sendToElastic -> ToRabbitMQ
- Display the version at the start
- Move file & fflag from internal to pkg

### 2.0.1 (2023/8/28)
Enhancements:
- Redirect the Gin logs to Graylog

### 2.0.0 (2023/8/24)
Description:
- Fully functionality for the **new agent**
- Redirect the log to the graylog

### 1.0.2 (2023/8/28)
Enhancements:
- Redirect the Gin logs to Graylog

### 1.0.1 (2023/8/24)
Enhancements:
- Redirecting the log to the graylog

### 1.0.0 (2023/8/17)
Description:
- Fully functionality for the **old agent**