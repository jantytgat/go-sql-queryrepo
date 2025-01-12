# go-sql-queryrepo - SQL Query Repository for Go

This library enables the use of centralized storage for all SQL queries used in an application.
You can either choose to load all queries into a repository, or load them from a filesystem as necessary.

The main intention is to embed a directory structure into the binary as ```embed.FS```, which can then used in the
application.

Two examples are available in the [example](https://github.com/jantytgat/go-sql-queryrepo/tree/main/example) directory:
- embbededSqlWithRepository
- embeddedSqlWithoutRepository

## Basics

### Add the package to your project

```bash
go get github.com/jantytgat/go-sql-queryrepo
```

### Import
Next, you can manually add the import statement to your ```.go```-file, or have it added automatically when using it.

```text
import github.com/jantytgat/go-sql-queryrepo
```

### Embed assets containing queries.

> [!IMPORTANT]  
> The root folder embedded in the application cannot have nested directories.  
> This means that the directory structure is limited to 1 level of collections, each containing a set of text files with
> a ```.sql``` extension.
>
> Files with another extension will fail to load!

Let's assume to following directory structure in an embedded filesystem:

```
/assets
|-- /queries
    |-- collection1
    |   |-- create.sql
    |   |-- read.sql
    |   |-- update.sql
    |   |-- delete.sql
    |-- collection2
        |-- list.sql
```

You can now embed the statement files as follows:

```go
//go:embed assets/queries/*
var f embed.FS
```

### With repository

#### Create a new repository

```go
// Create query repository from embedded files
var r *queryrepo.Repository
if r, err = queryrepo.NewFromFs(f, "assets/queries"); err != nil {
panic(err)
}
```

#### Load a query from the repository

Now the repository has been initialized, we can get a query from a collection:

```go
var query string
if query, err = r.Get("collection1", "create"); err != nil {
panic(err)
}
fmt.Println("Query:", query)
```

### Without repository

If you don't want to initialize a repository, but rather choose to load SQL queries straight from the filesystem, you
can do so as follows:

```go
var query string
if query, err = queryrepo.LoadQueryFromFs(f, "assets/queries", "collection2", "list"); err != nil {
panic(err)
}
```

## Prepared statements

We also provide the means to create prepared statements for the queries, either with or without using a repository, as
long as a ```Preparer``` is passed into the functions.

```go
type Preparer interface {
Prepare(query string) (*sql.Stmt, error)
}
```
For example:

- *sql.Db
- *sql.Tx

### With repository

```go
var createStmt *sql.Stmt
if createStmt, err = queryrepo.Prepare(db, r, "collection1", "create"); err != nil {
panic(err)
}
```

### Without repository

```go
var createStmt *sql.Stmt
if createStmt, err = queryrepo.PrepareFromFs(db, f, "assets/statements", "collection1", "create"); err != nil {
panic(err)
}
```