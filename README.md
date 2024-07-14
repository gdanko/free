# Free
Free is a Golang-based port of the Linux [free](https://man7.org/linux/man-pages/man1/free.1.html) command. It displays the amount of free and used memory in the system.

## Installation
Clone the repository and from within the repository directory, type `make build`. This will create a directory with the given value of `GOOS` and install the binary there. It will also create a tarball which will eventually be used for Homebrew formulae.

## Free Runtime Options
```
$ free --help
Usage:
  free [OPTIONS]

Application Options:
  -b, --bytes    Show output in bytes.
      --kilo     Show output in kilobytes. Implies --si.
      --mega     Show output in megabytes. Implies --si.
      --giga     Show output in gigabytes. Implies --si.
      --tera     Show output in terabytes. Implies --si.
      --peta     Show output in petabytes. Implies --si.
      --exa      Show output in exabytes. Implies --si.
  -k, --kibi     Show output in kibibytes.
  -m, --mebi     Show output in mebibytes.
  -g, --gibi     Show output in gibibytes.
  -t, --tebi     Show output in tebibytes.
  -p, --pebi     Show output in pebibytes.
  -e, --exbi     Show output in exbibytes.
  -j, --json     Output the data as JSON.
  -y, --yaml     Output the data as YAML.
      --si       Use kilo, mega, giga, etc (power of 1000) instead of kibi, mebi, gibi (power of 1024).
  -l, --lohi     Show detailed low and high memory statistics (Linux Only).
      --total    Show total for RAM + swap.
  -s, --seconds= Continuously display the result N seconds apart.
  -c, --count=   Display the result N times.
  -w, --wide     Wide output (Linux Only).
  -V, --version  Output version information and exit.

Help Options:
  -h, --help     Show this help message
  ```
  
## Free Examples
### No args
```
$ free
                 total         used         free       active     inactive        wired
Mem:          67108864     62743840      3315056     30068288     29658656      3016896
Swap:                0            0            0
```

### JSON output
```
$ free --json
{
    "memory": {
        "total": 68719476736,
        "used": 64277970944,
        "free": 3365519360,
        "active": 30850154496,
        "inactive": 30482284544,
        "wired": 2945531904
    },
    "swap": {
        "total": 0,
        "used": 0,
        "free": 0
    }
}
```

### Display five samples at a five second inteveral
```
$ free --count 5 --seconds 5
                 total         used         free       active     inactive        wired
Mem:          67108864     62759792      3296928     30121488     29773760      2864544
Swap:                0            0            0

                 total         used         free       active     inactive        wired
Mem:          67108864     62784144      3273040     30090416     29808512      2885216
Swap:                0            0            0

                 total         used         free       active     inactive        wired
Mem:          67108864     62743520      3314688     30050432     29775232      2917856
Swap:                0            0            0

                 total         used         free       active     inactive        wired
Mem:          67108864     62724384      3332064     30081568     29767792      2875024
Swap:                0            0            0

                 total         used         free       active     inactive        wired
Mem:          67108864     62728864      3328816     30043360     29770352      2915152
Swap:                0            0            0
```
