# Remind

Remind is a project aware todo app that will show relevant todos depending on the project folder. It also stores all notes in $HOME/remind/ to allow for easy syncing

## Installation (MacOs only currently):

```sh
brew install danwlker/remind/remind@0.0.4-alpha
```

Alternatively, with Go installed:

```sh
go install github.com/DanWlker/remind@latest
```

Uninstallation (MacOs only currently):

```sh
brew uninstall danwlker/remind/remind@0.0.4-alpha
```

## Build

```sh
go build .
```

## Usage

1. Adding a reminder associated to this folder

   ```sh
   remind add <your reminder>
   ```

   To associate this to the global folder use the `-g` flag

   ```sh
   remind add -g <your reminder>
   ```
   
1. Listing todos associated with this folder

   ```sh
   remind list
   ```

   To list all reminders use the `-a` flag

   ```sh
   remind list -a
   ```

   To list the global folder only use the `-g` flag

   ```sh
   remind list -g
   ```

1. To edit your reminders in a text editor

   ```sh
   remind edit
   ```

   To edit your global reminder use the `-g` flag

   ```sh
   remind edit -g
   ```

1. To remove a particular reminder, use its number. To see its number use the `list` command

   ```sh
   remind remove <reminder number>
   ```

   To remove all the reminders associated with this folder, use the `-a` flag

   ```sh
   remind remove -a
   ```

   To target the global reminder folder use the `-g` flag

   ```sh
   remind remove -g <reminder number>
   ```
