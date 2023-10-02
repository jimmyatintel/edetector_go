#!/bin/bash

if [ $# -ne 2 ]; then
    echo "Usage: $0 <number_of_times> <server_ip>"
    exit 1
fi

count=$1
ip=$2

if ! [[ "$count" =~ ^[0-9]+$ ]] || [ "$count" -lt 1 ]; then
    echo "Please provide a valid positive integer as the first argument."
    exit 1
fi

output_file="cmd/mockagent/prcessID.txt"
> "$output_file"

echo "Server IP: $ip"
for ((i = 1; i <= count; i++)); do
    nohup go run cmd/mockagent/agent.go "$ip" >/dev/null 2>&1 &
    echo $! >> "$output_file"
done

echo "Generated $count mock agent. Process IDs written to $output_file."
