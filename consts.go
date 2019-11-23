package goDB

/*
 * Comments related to page structures
 */

// NULL_POINTER_PAGE is a page that doesnt exist
const NULL_POINTER_PAGE = 0

// TYPE_TABLE_LIST holds the id of a table_List page
const TYPE_TABLE_LIST = 0

// TYPE_TABLE_SCHEMA holds the id of a table_schema page
const TYPE_TABLE_SCHEMA = 1

// TYPE_EMPTY_PAGES_LIST holds the id of a empty_pages_list page
const TYPE_EMPTY_PAGES_LIST = 2

// TYPE_TABLE_ROWS holds the id of a table_rows page
const TYPE_TABLE_ROWS = 3

/*
 * Comments related to pageHead structures
 */

// PAGE_HEAD_SIZE represents the length of a pageHead (including place for iv)
const PAGE_HEAD_SIZE = 64

/*
 * Comments related to dbHead
 */

const DB_HEAD_SIZE = 128

/*
 * General constants
 */

// IV_SIZE is the size of the iv used for encryption
const IV_SIZE = 16

// MASTER_KEY_SIZE is the size of the master encryption key
const MASTER_KEY_SIZE = 32
