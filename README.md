# Downloading to files in parallel
If you have a file `urls.txt` containing one url per line, to download all of them, run:
```
go run . -f urls.txt -b 20
```
It will download from each of the urls, with up to 20 concurrent downloads.