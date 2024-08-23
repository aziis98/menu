# menu

<p align="center">
<img src="https://github.com/user-attachments/assets/455c9ec9-0375-40cc-a0ef-b6cfe31776cc" />    
</p>

A dmenu like tool for the terminal made with [BubbleTea](https://github.com/charmbracelet/bubbletea). You can pass a command to it and it will show you a prompt that continuously calls the command and shows you the output with selectable lines.

To easily integrate search functionality, there is also a `menu search <query>` subcommand that uses a [fuzzy search](https://github.com/sahilm/fuzzy) algorithm to filter stdin lines.

I really the concept of "command pallets" or "ctrl+p", so I decided to make my own. Rofi is great but is somewhat limited (it is specifically a program launcher with extra features). So I wanted to make something more generic and composable that I could easily integrate with other stuff using shell scripts. (There is also the preview feature of fzf but didn't exactly fit my needs).

## Usage

```bash
git clone https://github.com/aziis98/menu
cd menu
go build -v -o ./bin/menu .
cp ./bin/menu ~/.local/bin
```

```bash
menu -s -c 'ls -1 | menu search $prompt'
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

- Better help usage message.