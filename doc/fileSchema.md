# Table of content

* [General](#general)
* [Header](#header)
  * [encryption of header](#encryption-of-header)
* [Pages](#pages)
  * [Page Header](#page-header)
  * [encryption of pages](#encryption-of-pages)
  * [Types of pages](#types-of-pages)
    * [Table list page](#table-list-page)
    * [Table schema page](#table-schema-page)
    * [Empty pages page](#empty-pages-page)
    * [Table rows page](#table-rows-page)
* [Column Data Type](#column-data-type)
* [Column Flags](#column-flags)
* [Future](#future)

> # General
> The Database file is composed of two parts.
> 1. The leading 128 bytes which make up the [Header](#header) and contain basic information on how to parse the database.
> 2. Repeating blocks of a fixed size (defaults to 4096 bytes) called the [Pages](#pages). There are different types of pages defining which information they hold.

> # Header
the header contains the following information
>
> | start index | length (in bytes) | added in | description |
> | --- | --- | --- | --- |
> | 0 | 5 | v1 | the file description (being "goDB" followed by a trailing 00) |
> | 5 | 3 | v1 | Version in Format Major, Minor, Patch |
> | 8 | 2 | v1 | the size of a single database page, needs to be at least 10 (which represents 1024 byte page size) |
> | 10 | 4 | v1 | count of database pages |
> | 14 | 4 | v1 | page index of first table_list page |
> | 18 | 4 | v1 | page index of first empty_pages_list |
> | 22 | 50 | / | reserved for future use |
> | 72 | 56 | / | unusable to keep the header at 128bytes when using encryption |
>
> ## encryption of header
> When using encryption the beginning of the regular [Header](#header) gets changed up. It now includes a __new file description__, the __iv used to encrypted the master encryption key__ with the users password as well as the __encrypted master encryption key__. The prefix can be seen in more detail in the table down below. The encrypted header is also 128 bytes long.
>
> **❗ It should be noted that the only thing bein dropped from the original header is the old "goDB" file description. ❗**
>
> | start index | length (in bytes) | added in | description |
> | --- | --- | --- | --- |
> | 0 | 9 | v1 | the new file description (being "goDB enc" followed by trailing 00) |
> | 9 | 16 | v1 | the iv used to encrypt the master encryption key |
> | 25 | 32 | v1 | encrypted Master encryption key |
> | 57 | 4 | v1 | "true" encoded as utf8 to validate master encryption key was decrypted correct |
> | 61 | 72 | / | regular [Header](#header) but without the "goDB" file description |

> # Pages
> Like described above the Pages are the structure holding the data itself.
>
> ## Page Header
> Every Page starts with a 64 byte long page header. There is a general part which can be found in every page which may be followed by a page_type specific header.
>
> The general Header is composed of the following:
>
> | start index | length (in bytes) | added in | description |
> | --- | --- | --- | --- |
> | 0 | 1 | v1 | the type of the page [(indexes can be found below)](#types-of-pages) |
> | 1 | 4 | v1 | the id of the previous page |
> | 5 | 4 | v1 | the id of the next page |
> | 9 | 2 | v1 | the first relevant byte holding data (relevant when data is overlapping one page) |
> | 11 | 2 | v1 | number of bytes that are trimmed of the end to make deletion/editing possible |
> | 13 | 35 | / | reserved for page specific use |
> | 48 | 16 | / | unusable to keep the page header at 64 bytes when using encryption |
>
> ## Encryption of Pages
> Encryption of a page is not optional. If the file description is "goDB enc" all pages are encrypted.
>
> When encrypting a page the first 16 bytes are sacrificed for the iv. Like with the file header the original data gets shifted back to account fir this.  
> **A new iv gets used every time the page gets written to the file.** Besides this the structure of the page stays the same.
>
> ## Types of pages
> At the moment 4 types of pages are implemented:
>
> * [table_list](#table-list-page) **(type=0)** to hold a list of all database tables. It's first page is referenced in the [Header](#header).
> * [table_schema](#table-schema-page) **(type=1)** to hold the information about the columns of a table. It's first page is referenced in the [table_list](#table-list-page).
> * [empty_pages_list](#empty-pages-page) **(type=2)** to hold a list of empty pages that became free to use after deletion. It's first page is also referenced in the [Header](#header).
> * [table_rows](#table-rows-page) **(type=3)** to hold the data of a table. It's first page is referenced in the [table_list](#table-list).
>
> ### Table list page
> The table_list is a page with **type=0** that holds information on which tables exist.
>
> **table_list extends the header with a 4 byte "table_count" variable telling us how many tables to expect to find**
>
> Content, repeated for every table:
>
> | start index | length (in bytes) | added in | description |
> | --- | --- | --- | --- |
> | 0 | 4 | v1 | entry length |
> | 4 | 4 | v1 | table uid |
> | 8 | 4 | v1 | page index of first table_schema page |
> | 12 | 4 | v1 | page index of first table_rows page |
> | 16 | 4 | v1 | page index of last table_rows page |
> | 20 | 4 | v1 | row count |
> | 24 | 4 | v1 | column count |
> | 28 | 2 | v1 | table name length |
> | 30 | n | v1 | table name (utf8) |
>
> ### Table schema page
> The table_schema is a page with **type=1** that holds information on how a single table is build
>
> **table_schema extends the header with a 4 byte "table_uid" variable referencing the data in the table_list page**
>
> Content, repeated for every column:
>
> | start index | length (in bytes) | added in | description |
> | --- | --- | --- | --- |
> | 0 | 4 | v1 | entry length |
> | 4 | 4 | v1 | column uid |
> | 8 | 1 | v1 | columns [data type](#column-data-type) |
> | 9 | 1 | v1 | columns [flags](#column-flags) |
> | 10 | 4 | v1 | AUTO_INCREMENT value for the given column |
> | 14 | 2 | v1 | column name size |
> | 16 | n | v1 | column name (utf8) |
>
> ### Empty pages page
> The empty_Pages is a page with **type=2** that holds information on tables that where free'ed up.
>
> The content is a simple list of 4 byte long page indexes. No other informations are provided / saved.
>
> ### Table rows page
> The table_rows is a page with **type=3** that holds the data belonging inside the table / column structure.
>
> **table_rows extends the header with a 4 byte "table_uid" variable referencing the data in the table_list page**
>
> Content, for each row of Data:
>
> | start index | length (in bytes) | added in | description |
> | --- | --- | --- | --- |
> | 0 | 4 | v1 | entry length |
> | 4 | 4 | v1 | row uid |
> | 8 | n | v1 | concat of all column_data |
>
> The **column_data** is a very basic construct of:
> * 1 optional isNull byte (depends on the nullable column flag)
> * 4 optional content_length bytes (depends on the column data type)
> * the content bytes itself

> # Column Data Type
> | Value | added in | Type |
> | --- | --- | --- |
> | 0 | v1 | UINT_8 |
> | 1 | v1 | UINT_16 |
> | 2 | v1 | UINT_32 |
> | 3 | v1 | UINT_64 |
> | 4 | v1 | INT_8 |
> | 5 | v1 | INT_16 |
> | 6 | v1 | INT_32 |
> | 7 | v1 | INT_64 |
> | 8 | v1 | DATE |
> | 9 | v1 | UTF8_STRING |
> | 10 | v1 | BOOLEAN |
> | 10-255 | / | reserved for future use |

> # Column Flags
> | Bit | added in | Type |
> | --- | --- | --- |
> | 0 | v1 | UNIQUE |
> | 1 | v1 | NOT_NULL |
> | 2 | v1 | AUTO_INCREMENT |
> | 3 | v1 | PRIMARY_KEY |
> | 4-7 | / | reserved for future use |

> # Future
> * should support blob as [data type](#column-data-type)
> * should support some form of float/number/real as [data type](#column-data-type)
> * might add a "file change counter" field to header
> * might add a "default page cache size" field to header
> * should look into for thread-safety / multi threading capability
> * atm when data is deleted it has to read the next page to move some data forward - not sure whether that's the best way 🤔
