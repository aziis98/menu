#!/bin/bash

TMP_FILE="/tmp/search.txt"

case $event in
    "open")
        find . -type f > "$TMP_FILE"
        cat "$TMP_FILE"
        ;;
    "close")
        cat "$TMP_FILE" | ./bin/menu search "$prompt" 2>/dev/null | sed -n "${selected}p"
        rm "$TMP_FILE"
        ;;
    *)
        cat "$TMP_FILE" | ./bin/menu search "$prompt" 2>/dev/null
        
        if [ $? -ne 0 ]; then
            echo "No match found"
        fi
        ;;
esac
