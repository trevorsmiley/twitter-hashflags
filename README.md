# twitter-hashflags
Twitter Hashflag Downloader is a command line tool that can list, sync and download all active hashflags from Twitter.

Shoutout to [@JamieMagee](https://github.com/JamieMagee) and his [HashflagArchive Twitter Bot](https://github.com/JamieMagee/hashflags-function) for providing the inspiration to create this tool.

# Usage

## Info

### Get hashflag information list

Return a list of active hashflags
```bash
./twitter-hashflags info
```

### Write hashflag information list with full details to file

Write a list of active hashflags with full details to `hashflag-list.txt`
```bash
./twitter-hashflags info -out
```

### Get missing hashflags

Sync hashflags with the `/downloaded_hashflags` folder and list all that are missing
```bash
./twitter-hashflags info -diff
```

### Get deactivated hashflags

Sync hashflags with the `downloaded_hashflags` folder and list all that are now deactivated
```bash
./twitter-hashflags info -deactivated
```

## Download

### Sync hashflags

Sync hashflags with the `/downloaded_hashflags` folder and download any that are missing
```bash
./twitter-hashflags sync
```

### Sync hashflags and move deactivated

Sync hashflags with the `/downloaded_hashflags` folder, download any that are missing and move all deactivated hashflags to `/deactivated`
```bash
./twitter-hashflags sync -m
```

### Force download

Clear the  `/downloaded_hashflags` folder and download all hashflags
```bash
./twitter-hashflags sync -force
```
