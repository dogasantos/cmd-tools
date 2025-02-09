# uniquelines.go

Show lines that are present on file 1 but are not present in file 2 (line order does not matter)

# letters.go

Remove everything but letters (a-z) from strings (use with pipe)
ex:
```cat text.txt |letters ```

# noschars.go
remove special chars from strings (use with pipe)

# numbers.go

display only numbers in strings (use with pipe)

ex:

```
$ cat file.txt
A331
this is a numeber: 44
nothing
more nothing @

$ cat file.txt|numbers
331
44
```
# filterurl

this utility will help to show or filter out ip-based urls from a file. 

```
https://1.1.1.1:9191
https://example.com:9191
```

-omit-ip will display only https://example.com:9191
-only-ip will display only https://1.1.1.1:9191




# tolower.go
lower everything (pipe)

# toupper.go
upper everything (pipe)

# filterlinesab.go
read two files and remove from the first file any content present in the second file (non destructive, it just print the result to stdout)

# filterips.go
take text as input and prints out just found ip addressess and/or ports along with ips:

# filterlines.go
take a file from -l <file> or stdin and remove lines based on a subset of rules (fuzzy match, regex match, string match, char count and randomness)
