# go-menu

A dmenu like tool for the terminal made with [BubbleTea](https://github.com/charmbracelet/bubbletea). You can pass a command to it and it will show you a prompt that continuously calls the command and shows you the output with selectable lines.

To easily integrate search functionality, there is also a `menu search <query>` subcommand that uses a [fuzzy search](https://github.com/sahilm/fuzzy) algorithm to filter stdin lines.

## Usage

```bash
$ go build -v -o ./bin/menu .
$ ./bin/menu -s -c 'ls -1 | ./bin/menu search $prompt'
```

### Commands

- `Enter`: Confirm selection and close the tool.

- `Up/Down Arrow`: Navigate through the command output.

- `Ctrl+C/Esc`: Quit the application without making a selection.

## Examples

For more examples, check the [examples](./examples) directory.

### Calculator

```bash
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
```

### Fuzzy Search

```bash
TMP_FILE="/tmp/search.txt"

case $event in
    "open")
        find . -type f > "$TMP_FILE"
        
        cat "$TMP_FILE"
        ;;
    "close")
        cat "$TMP_FILE" | ./bin/menu search "$prompt" | sed -n "${selected}p"
        
        rm "$TMP_FILE"
        ;;
    *)
        cat "$TMP_FILE" | ./bin/menu search "$prompt"
        
        if [ $? -ne 0 ]; then
            echo "No match found"
        fi
        ;;
esac
```

## ToDo

- Organize the code better.

- Add more examples.

- Choose a better name for the project.
