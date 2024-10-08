# menu

<p align="center">
<img width="720" src="https://github.com/user-attachments/assets/d1a7021d-bb5f-4476-b20b-aca3567ead89" />    
</p>

A dmenu like tool for the terminal made with [BubbleTea](https://github.com/charmbracelet/bubbletea). You can pass a command to it and it will show you a prompt that continuously calls the command and shows you the output with selectable lines.

To easily integrate search functionality, there is also a `menu search <query>` subcommand that uses a [fuzzy search](https://github.com/sahilm/fuzzy) algorithm to filter stdin lines.

I really the concept of "command pallets" or "<kbd>Ctrl</kbd> + <kbd>P</kbd>" shortcut, so I decided to make my own generic tool for making multiple of this. Rofi is great but is somewhat limited (it is specifically a program launcher with extra features). So I wanted to make something more modular and composable that I could easily integrate with other stuff using shell scripts. (There is also the preview feature of fzf but even that didn't exactly fit my needs).

## Usage

This is how to currently get and build the tool

```bash
$ git clone https://github.com/aziis98/menu
$ cd menu
$ go build -v -o ./bin/menu .

# Copy the binary to your path
$ cp -f ./bin/menu ~/.local/bin/menu

# Or symlink it
$ ln -sfr ./bin/menu ~/.local/bin/menu
```

And here is an example

```bash
menu -s -c 'ls -1 | menu search $prompt'
```

For completion this is the help message

```
Usage:
  menu [OPTIONS] COMMAND
  menu search QUERY

  This is a simple interactive command line tool that allows you to run a
  command at every key press. It can be used to create interactive prompts or
  menus.

  The search subcommand performs a fuzzy search on the stdin lines and prints
  the results highlighting the matched characters. If no query is provided, it
  will print all the lines.

Options:
  -i, --initial string       Initial prompt text
  -p, --placeholder string   Placeholder text
  -s, --selection            Only return the selected item

Command:
  This command will be executed at every key press so mind the performance. It
  has access to the following environment variables:

    $prompt     the current prompt text
    $event      the event that triggered the command (open, key, select, close)
    $sel_line   the selected line
    $sel_index  the index of the selected item, starting from 1
```

### Commands

- <kbd>Enter</kbd>: Confirm selection and close the tool.

- <kbd>Up Arrow</kbd> / <kbd>Down Arrow</kbd>: Navigate through the command output.

- <kbd>Ctrl</kbd> + <kbd>C</kbd> / <kbd>Esc</kbd>: Quit the application without making a selection.

## Examples

For more examples, check the [examples](./examples) directory, for now there is a simple calculator and a search example.

## ToDo

- Organize the code better.

- Add more examples.

- Choose a better name for the project.

- Better help usage message.

- Put on the AUR.
