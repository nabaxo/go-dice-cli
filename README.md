# Dead simple dice-rolling
Just a small CLI app I wrote in Go for fun, for my own use.

  > [!BUG] This is the first thing I wrote in Go, I know the code is trash.

It uses [Random.org](random.org), so if you wanna run it yourself, you need to create a file that's called `apikey` with, well, your [Random.org API key](https://api.random.org/dashboard) and put it in the root/next to the binary.
### Build script for a smaller binary (requires UPX to be installed)
```sh
./build.sh
```
## Instructions of use:
### Run it in the terminal like this
```sh
./dice-roller
```

### Format your dice rolls like this:
- `ndp+q`, for example `d10`, `2d20`, `2d12-5`, or `10d6a3` etc
- `n`:  n is number of dice (optional)
- `p`:  is type of dice
- `+`:  is the type of modifier
- `q`:  is the modifier

#### The types of modifier are:
- `+`: Add a constant number to roll
- `-`: Subtract a constant number from roll
- `a`: Show all dice which rolled the modifier number and above
- `b`: Show all dice which rolled the modifier number and below

### Other commands
- `r or [enter]`:  for repeat last
- `r[n]`:  where n is a specific roll to repeat
- `q`:  to quit
