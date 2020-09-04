# backup and restore your DNS entries

I need to migrate one of my DNS nameservers from digital ocean to gandi, so I wrote this...

```
./dnsbackup backup ona.im > ona.im.json
```

and then can run the following to push those same values to the alho.st domain

```
./dnsbackup restore --dry-run alho.st ona.im.json
```