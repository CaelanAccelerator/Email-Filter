#!/bin/bash
# demo.sh — send test emails to the server
for f in testdata/*.eml; do
    echo ">>> sending $f"
    nc -N localhost 8080 < "$f"
    sleep 1          
done