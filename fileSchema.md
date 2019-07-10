# Todo
Header > indizes of table_schema

# Header
the first 128 bytes of the database file, including basic information (and encryption details)
> ## non-encrypted Header
> | Byte Offset | Length | Meaning |
> | --- | --- | --- |
> | 0 | 5 | "goDB" (and trailing 00) as utf8 string |
> | 5 | 3 | Version in Format Major, Minor, Path |
> | 8 | 2 | database page size in bytes (power of two between 512 and 32768)
> | 10 | 4 | File change counter |
> | 14 | 4 | Count of Database Pages |
> | 18 | 4 | Default page cache size |
> | 19 | 108 | Reserved for future use |
>
> ## encrypted header
> | Byte Offset | Length | Encrypted | Meaning |
> | --- | --- | --- | --- |
> | 0 | 16 | false | the iv for the header
> | 16 | 9 | true | "goDB enc" (and trailing 00) as utf8 string |
> | 25 | 32 | true | Master encryption key |
> | 57 | 3 | true | Version in Format Major, Minor, Path |
> | 60 | 2 | true | database page size in bytes (power of two between 512 and 32768)
> | 62 | 4 | true | File change counter |
> | 66 | 4 | true | Count of Database Pages |
> | 70 | 4 | true | Default page cache size |
> | 74 | 54 | true | Reserved for future use |
