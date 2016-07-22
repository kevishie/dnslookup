# dnslookup

Lookup tool with the ability to query multiple servers.

## Example usage

```bash
$ dnslookup espn.com
google-public-dns-a.google.com ==> [199.181.132.250]
google-public-dns-b.google.com ==> [199.181.132.250]
```

## Configuration

The program will look for a .resolv.conf file in your home directory for a list of dns servers to query. The list must be in the following format:

```bash
google-public-dns-a.google.com 8.8.8.8
google-public-dns-b.google.com 8.8.4.4
```

If not .resolv.conf file is found in the user's home directory, google dns a and b will be used by default.
