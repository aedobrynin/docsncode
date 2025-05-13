# Documentation

## Supported Languages

[List of supported languages](supported_languages.md)

## Comment Blocks

Comment blocks should be a valid Markdown code.
All [CommonMark Markdown Spec](https://spec.commonmark.org/)
features are available, except code blocks. For example,
you can make your text bold or italic, add an image, or a 
hyperlink. There is also other features that will be listed
further.

To add a comment block you should use regular comments mechanism
from your programming language. There should be `@docsncode`
text to mark the beginning and the end of the comment block.
For example, if you're writing on C++, you could write:
```
// @doscncode
// This is a comment block
// @docsncode
```
It's also possible to use multiline comments:
```
/* @docsncode
This is a multiline comment block.

Blah-blah-blah.
@docsncode */
```
The `@doscncode` mark should be placed at the first and the last
lines of the comment. For example, this is not allowed:
```
/*
@docsncode
...
@docsncode
/*
```

## Code Blocks

Code block is everything that's not a comment block. The resulted
code block will have syntax highlight provided by
[highlight.js](https://highlightjs.org/).

## Hyperlinks

You can add hyperlinks in your comment blocks. If it's a link to 
file from your project, it will be automatically transformed
to correct link in a result file. For example, if you have main.go
and it refers to sum.go in the same folder, you can simply write:
```
// main.go

// @docsncode
// [link](sum.go)
// @docsncode
```
and that link to `sum.go` will be automatically transformed to
`sum.go.html`.

**Absolute hyperlinks are not allowed.**

If the referred file won't have result file
(e.g. it won't have corresponding file in the result directory), 
the hyperlink will refer to the original file. So it's important
to have access to original project when you're watching the
docsncode output. Otherwise, some hyperlinks won't work.

## Diagrams

It's possible to add diagrams to your comment blocks. It should
be [Mermaid diagrams](https://mermaid.js.org/intro/)
The syntax is the following:
```
// @docsncode
// ```mermaid
// <your diagram>
// ```
// @docsncode
```

## Cache

By default, DocsnCode results are cached. That is, if you change
only one file from your project, the result will be regenerated
only for that file. Default caching strategy is `modtime` that
checks that source and result file modification time timestamps
didn't change from the last run. There is also cache based on
SHA-256 hashes, to use it simply provide `--cache hash`. if you
want to disable caching, you can write `--cache none`. To force
rebuild the result write `--force-rebuild`.

The cache data is stored in `.docsncode_cache.json` file at the 
root of the project. If you want to change that behaviour, you
can provide the path to cache data file as third positional 
argument. For example, `./docsncode project result cache.json`.

## Ignoring some files

There is an ability to do not generate any output for specific
files or directories. You can use `.docsncodeignore` file at the
root of your project. It's the same file as a regular `.gitignore`
file.
