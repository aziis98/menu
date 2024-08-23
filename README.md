# menu

<p align="center">
<img src="https://github.com/user-attachments/assets/455c9ec9-0375-40cc-a0ef-b6cfe31776cc" />    
</p>

A dmenu like tool for the terminal made with [BubbleTea](https://github.com/charmbracelet/bubbletea). You can pass a command to it and it will show you a prompt that continuously calls the command and shows you the output with selectable lines.

To easily integrate search functionality, there is also a `menu search <query>` subcommand that uses a [fuzzy search](https://github.com/sahilm/fuzzy) algorithm to filter stdin lines.

I really the concept of "command pallets" or "ctrl+p", so I decided to make my own. Rofi is great but is somewhat limited (it is specifically a program launcher with extra features). So I wanted to make something more generic and composable that I could easily integrate with other stuff using shell scripts. (There is also the preview feature of fzf but didn't exactly fit my needs).

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

- `Enter`: Confirm selection and close the tool.

- `Up/Down Arrow`: Navigate through the command output.

- `Ctrl+C/Esc`: Quit the application without making a selection.

## Examples

For more examples, check the [examples](./examples) directory, for now there is a simple calculator and a search example.

## ToDo

- Organize the code better.

- Add more examples.

- Choose a better name for the project.

- Better help usage message.

- Put on the AUR.