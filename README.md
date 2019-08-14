# ObjectiveSeeDownloader

After reading the fantastic ebook on reversing macOS malware I thought I would look for some samples and quickly came across the collection Objective See has. 
This is great but there is no batch download and its not a open directory, so I thought I would write a littel script for it. If the outputdir flag isn't set then it will default download 
to `./malware/`

## Usage 

```
cd objectiveSeeDownloader
go build 
./objectiveSeeDownloader
```

After this you should have all the samples indexed by Objective See