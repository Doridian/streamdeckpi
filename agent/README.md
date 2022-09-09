# StreamDeckPi agent

## Basics

The configuration is split up into **page**s.

Each **page** has a single **action** assigned to each button.

On startup, the page `default.yml` will be loaded.

Any file loaded will first try being loaded relative to the config directory, and then from the embedded FS.

You can adjust the config directory by setting `STREAMDECKPI_CONFIG_DIR` environment variable.

If this variable is unset or blank, it will default to `./config`.

**Page**s are stored in a stack to allow for operations like "go back".

## Configuration

The configuration uses YAML. Each page has its own, independent, YAML file.

These files look as follows (buttons are organized however the Stream Deck device orders them):

```yaml
timeout: 10s # Optional, if set: If no button is pressed in this period, go back to previous page
actions:
- name: some_action
  button: [0,0] # Button coordinates start at [0,0] for top-left
                # [0,1] is the button below [0,0]
                # [1,0] is the button to the right of [0,0]
  parameters:
     name: "value"
- name: some_other_action
  button: [3,0]
  parameters:
     number: 42
```

## Actions

### Page swaps

#### swap_page

Replaces the current **page** on the stack with the given **page**.

```yaml
name: swap_page
button: [0,0]
parameters:
    icon: some_icon.png
    target: some_page.yml
```

#### push_page

Pushes the given **page** onto the stack.

```yaml
name: push_page
button: [0,0]
parameters:
    icon: some_icon.png
    target: some_page.yml
```

#### pop_page

Pops the current **page** from the stack.

```yaml
name: pop_page
button: [0,0]
parameters:
    icon: some_icon.png
```

### Miscellaneous

#### none

Does nothing at all, useful if you want a button to do nothing.

```yaml
name: none
button: [0,0]
parameters:
    icon: some_icon.png
```

#### exit

Causes the agent process to exit gracefully.

```yaml
name: exit
button: [0,0]
parameters:
    icon: some_icon.png
    exit_code: 0 # Optional, default is 0
```

#### reset

Resets agent app to the default **page** causi.ng everything to go back to a state as if the app just started

```yaml
name: reset
button: [0,0]
parameters:
    icon: some_icon.png
```

#### shell

Runs a shell command and can optionally define icons to be set on specific states.

```yaml
name: command
button: [0,0]
parameters:
    command: [echo, "Hello world"] # Command and arguments as an array
    icon: some_icon.png # Icon to use when the command has never been run or as
                        # the default if no more specific icon is defined (see below)

    running_icon: some_running_icon.png # Icon to use while the command is running (optional)
  
    exit_code_icons: # Optional, if not given will revert back to "icon" immediately after exit
      0: some_success_icon.png # Icon to use if the command exits with code 0 (usually, this means success)
      1: some_exit1_icon.png # Icon to use for exit code 1 (you can use any number exit code to handle)
      default: some_error_icon.png # Icon to use if the command exits and no explicit icon is defined
    
    exit_to_idle_time: 5s # If this is set, will revert to "icon" after this time after the process exited
```
