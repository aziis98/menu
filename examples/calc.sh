#!/bin/bash

CALC_HISTORY_FILE="$HOME/.cache/calc_history"

view() {
    echo "= $1"
    echo "[Clear History]"
    echo
    echo "History:"
    tac "$CALC_HISTORY_FILE"
}

case $event in
    "open")
        # Create the file if it doesn't exist
        touch "$CALC_HISTORY_FILE"

        # Remove empty lines
        sed -rie '/^\s*=?\s*$/d' "$CALC_HISTORY_FILE"
        
        view
        ;;
    "close")
        result=$(echo "$prompt" | bc 2>/dev/null)
        
        if [ $? -eq 0 ]; then
            echo "$prompt = $result" >> "$CALC_HISTORY_FILE"

            # if sel_line has a "="
            if [[ $sel_line == *"="* ]]; then
                # Show the result
                echo $sel_line | sed 's/^.*=\s*//'
            elif [[ $sel_line == "[Clear History]" ]]; then
                # Clear the history
                > "$CALC_HISTORY_FILE"
            fi
        else
            echo "Invalid input"
        fi
        ;;
    *)
        result=$(echo "$prompt" | bc 2>/dev/null)
        
        if [ $? -eq 0 ]; then
            view "$result"
        else
            echo "Invalid input"
        fi
        ;;
esac
