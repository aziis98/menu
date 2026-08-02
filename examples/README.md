# Examples

To run the examples just call `menu ./examples/<name>.sh` from the root of the repository.

- [Calculator](./calc.sh)

    A simple calculator that uses `bc` to perform calculations. It also saves the history of calculations (like rofi).

- [Search](./search.sh)

    Small example that shows the usage of the `menu search` subcommand.

- [Launcher](./launcher.py)

    A simple launcher that uses `menu` to search for applications and run them. The discovery of applications is done with `Gio` in python.

- [Zoxide Recent](./zoxide-recent.sh)

    A directory picker that uses `zoxide` to list recently visited directories, `fzf` for fuzzy filtering, and `prefixtree` to display the paths as a tree.