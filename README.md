Gold
====

gold | blog -- Blogs revert

Inspired by [Blogsum](https://github.com/obfuscurity/blogsum)

[Documentation](https://godoc.org/github.com/dim13/gold)

Internal
========

URL schema
----------

URI			| Description			| Request
---			| -----------			| -------
/admin			| admin interface		| GET
/admin/add		| add article form		| GET
/admin/add		| add article			| POST
/admin/{slug}		| edit article form		| GET
/admin/{slug}		| enable article		| POST publish
/admin/{slug}		| disable article		| POST suppress
/admin/commtents	| show comments			| GET
/admin/commtents	| publish comment		| POST publish
/admin/commtents	| delete comment		| POST suppress
/{slug}			| show single article		| GET
/{slug}			| show single article		| POST comment
/{year}/{month}/{slug}	| show single article		| GET
/{year}/{month}/{slug}	| show single article		| POST comment
/{year}/{month}		| show all articles of month	| GET
/{year}			| show all artilces of year	| GET
/tag/{tag}		| show all articles with tag	| GET
/page/{number}		| show page number		| GET
/			| show main page		| GET
