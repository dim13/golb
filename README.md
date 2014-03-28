Golb
====

Blogs revert

Inspired by [Blogsum](https://github.com/obfuscurity/blogsum)

[Documentation](https://godoc.org/github.com/dim13/golb)

Internal
========

URL schema
----------

URI				| Description
---				| -----------
/admin				| admin interface
/admin/add			| add new article
/admin/{slug}			| edit article
/admin/{slug}/publish		| enable article
/admin/{slug}/suppress		| disable article
/admin/{slug}/commtens		| show articles comments
/admin/{slug}/commtent/publish	| publish comment
/admin/{slug}/commtent/suppress	| delete comment
/{slug}				| show single article
/{year}/{month}/{slug}		| show single article
/{year}/{month}			| show all articles of month
/{year}				| show all artilces of year
/tag/{tag}			| show all articles with tag
/page/{number}			| show page number
/				| show main page
