# twitter-hashflags
Twitter Hashflag Downloader is a command line tool that can list, sync and download all active hashflags from Twitter.

Shoutout to [@JamieMagee](https://github.com/JamieMagee) and his [HashflagArchive Twitter Bot](https://github.com/JamieMagee/hashflags-function) for providing the inspiration to create this tool.

# Usage

## List

### Get hashflag list

Return a list of active hashflags
```bash
./twitter-hashflags list
```

### Get hashflag list with full details

Return a list of active hashflags with full details
```bash
./twitter-hashflags list-fulldetails
```

### List missing hashflags

Sync hashflags with the `downloaded_hashflags` folder and list all that are missing
```bash
./twitter-hashflags diff
```

## Download

### Sync hashflags

Sync hashflags with the `downloaded_hashflags` folder and download any that are missing
```bash
./twitter-hashflags sync
```

### Force download

Clear the  `downloaded_hashflags` folder and download all hashflags
```bash
./twitter-hashflags force-download
```
