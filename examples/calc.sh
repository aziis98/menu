#!/bin/bash

CALC_HISTORY_FILE="$HOME/.cache/calc_history"

case $event in
    "open")
        # Create the file if it doesn't exist
        touch "$CALC_HISTORY_FILE"

        # Remove empty lines
        sed -i '/^$/d' "$CALC_HISTORY_FILE"
        
        # Show the history
        echo " = "
        tac "$CALC_HISTORY_FILE"
        ;;
    "close")
        result=$(echo "$prompt" | bc 2>/dev/null)
        
        if [ $? -eq 0 ]; then
            echo "$prompt = $result" >> "$CALC_HISTORY_FILE"

            # Print the selected line
            tac "$CALC_HISTORY_FILE" | sed -n "${selected}p" | sed 's/^.*=\s*//'
        else
            echo "Invalid input"
        fi
        ;;
    *)
        result=$(echo "$prompt" | bc 2>/dev/null)
        
        if [ $? -eq 0 ]; then
            echo " = $result"
        else
            echo "Invalid input"
        fi

        tac "$CALC_HISTORY_FILE"
        ;;
esac
