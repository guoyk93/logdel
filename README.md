# logdel

log file deletiong tool

## Usage

1. create rule files at `/etc/logdel.d`

```
# file patterns: keep days
/home/logs/info*.log: 3
/home/logs/error*.log: 5
```

2. execute `logdel`

## Credits

Guo Y.K., MIT License
