# Todo
* Header > indizes of table_schema
* Header > indizes of empty_pages_list

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

# Pages
## table_list (id=0)
> a page (often the first) that is referenced in the header and lists all the tables that exist:
> ### table_list>head
> | Byte Offset | Length | Meaning |
> | --- | --- | --- |
> | 0 | 4 | page type (always 0) |
> | 4 | 4 | index of previous page |
> | 8 | 4 | index of next page |
> | 12 | 4 | count of tables |
> | 16 | 48 | reserved for future use |
>
> ### table_list>body
> for each table in the list:
>
> | Byte Offset | Length | Meaning |
> | --- | --- | --- |
> | k | 4 | index of table |
> | k+4 | 4 | index of next table entry |
> | k+8 | 4 | page index of table_schema page |
> | k+12 | 4 | page index of first content page |
> | k+16 | 4 | row count |
> | k+20 | 4 | length of table name |
> | k+24 | n | table name (as utf8) |
## table_schema (id=1)
> a page (often the first) that is referenced in the header and list all the tables that exist:
> ### table_schema>head
> | Byte Offset | Length | Meaning |
> | --- | --- | --- |
> | 0 | 4 | page type (always 1) |
> | 4 | 4 | index of previous page |
> | 8 | 4 | index of next page |
> | 12 | 4 | count of columns |
> | 16 | 48 | reserved for future use |
>
> ### table_schema>body
> for each table in the list:
>
> | Byte Offset | Length | Meaning |
> | --- | --- | --- |
> | k | 4 | index of column |
> | k+4 | 4 | index of next column |
> | k+8 | 1 | columns data type (check [#column_type](#column_type) for more information) |
> | k+9 | 1 | column flags [UNIQUE, NOT_NULL, AUTO_INCREMENT, PRIMARY_KEY, reserved, reserved, reserved, reserved] |
> | k+10 | 4 | column name size |
> | k+14 | n | column name (as utf8) |
## empty_pages_list (id=2)
> a page that lists all empty / unused pages
> ### empty_pages_list>head
> | Byte Offset | Length | Meaning |
> | --- | --- | --- |
> | 0 | 4 | page type (always 2) |
> | 4 | 4 | index of previous page |
> | 8 | 4 | index of next page |
> | 12 | 52 | reserved for future use |
>
> ### empty_pages_list>body
> for each page in the list:
>
> | Byte Offset | Length | Meaning |
> | --- | --- | --- |
> | k | 4 | index of empty page |
## content (id=3)
> a page with content of a table
> ### content>head
> | Byte Offset | Length | Meaning |
> | --- | --- | --- |
> | 0 | 4 | page type (always 3) |
> | 4 | 4 | index of previous page |
> | 8 | 4 | index of next page |
> | 12 | 4 | index of the first row (if starting with data that didn't fit previous page) |
> | 16 | 48 | reserved for future use |
>
> ### content>body
> start of a row
>
> | Byte Offset | Length | Meaning |
> | --- | --- | --- |
> | k | 4 | index of row |
> | k+4 | 4 | index of next row |
>
> for each column:
>
> | Byte Offset | Length | Meaning |
> | --- | --- | --- |
> | 0 | 1 | bool: isNull |
> | 1 | 4 | optional: content_length |
> | 5 | n | content |

# Methods
* SELECT
* INSERT
* UPDATE
* DELETE

# column_type
> | Value | Type |
> | --- | --- |
> | 0 | UINT_8 |
> | 1 | UINT_16 |
> | 2 | UINT_32 |
> | 3 | UINT_64 |
> | 4 | INT_8 |
> | 5 | INT_16 |
> | 6 | INT_32 |
> | 7 | INT_64 |
> | 8 | DATE |
> | 9 | UTF8_STRING |
> | 10 | BOOLEAN |
> | 11 | FLOAT/NUMBER/REAL |
> | 12 | BLOB |
