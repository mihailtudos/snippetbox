# Snippetbox


## Description
My first go web app project

## Running the application

### Locally
To start the application locally run: 

```sh
    go run ./cmd/web
```

The above command would start the web application at a default port with the default config, if you want to give you own config run the app with the flag -help to see the available config:

```sh 
    go run ./cmd/web/ -help
```

## Static file serving

To serve static file the http.FileServe handler was used. It has multiple advantages:

1. It sanitizes all requests to stop directory traversal attacks
2. Range requests are fully supported (support resumable downloads)
3. The Last-Modified and If-Modified-Since headers are transparently supported
4. The Content-Type is automatically set from the file extension u

### Adding single file serving capabilities 

To serve single file a new handler can be set as (warning: http.ServeFile() does not automatically sanitize the file path - you must sanitize the input with filepath.Clean() before using it): 

```go
func downloadHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./ui/static/file.zip")
}
```

### Disable static file 

Disable subdirectory trans by creating index.html files in each static subdirectory:

```sh
    find ./ui/static -type d -exec touch {}/index.html \;
```

A better but more complex approach is to create a custom implementation of http.FileSystem, and have it return an os.ErrNotExist. 

## Logging - structured logger

To decouple the logging from the application a basic structured, concurrency-safe, logger was created with the default slog package (set to: os.Stdout). The final destination of the logs can be managed by your execution environment independently of the application: 

```sh 
    go run ./cmd/web >> /tmp/web.log
```

## DI

This web application, as most web applications, has multiple dependencies that the handlers need to access, such as a database connection pool, centralized error handlers, and template caches. Thus, to make these dependencies available, following the good practices they will be inject. 

DI is build around the application struct and works well with the handlers as long as they are in the same package **main**, for a more complex project structure [Closures for dependency injection](https://gist.github.com/alexedwards/5cd712192b4831058b21) can be used. 

## Database

To create a database run the following commands: ( a migration will be added later )

```mysql
    CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

    USE snippetbox;

    CREATE TABLE snippets (
        id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
        title VARCHAR(100) NOT NULL,
        content TEXT NOT NULL,
        created DATETIME NOT  NULL,
        expires DATETIME NOT NULL
    );

    CREATE INDEX idx_snippets_created ON snippets(created);
    
    INSERT INTO snippets (title, content, created, expires) VALUES (
        'An old silent pond',
        'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
        UTC_TIMESTAMP(),
        DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
    );

    INSERT INTO snippets (title, content, created, expires) VALUES (
        'Over the wintry forest',
        'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
        UTC_TIMESTAMP(),
        DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
    );

    INSERT INTO snippets (title, content, created, expires) VALUES (
        'First autumn morning',
        'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
        UTC_TIMESTAMP(),
        DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
    );

```

**Important**: Make sure to swap 'pass' with a password of your own choosing.

```mysql
    CREATE USER 'web'@'localhost';
    GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'localhost';
    ALTER USER 'web'@'localhost' IDENTIFIED BY 'pass';
```

Exit, and authenticate as the new user and use as default the snippetbox DB:

```mysql
    mysql -D snippetbox -u web -p
```
Get the MySQL driver for go:

```sh
    go get github.com/go-sql-driver/mysql@v1
```

This will update the **.mod** file with the required package dependency but also update the **.sum** file which is a cryptographic checksums representing the content of the required packages:

## [Reproducible build](https://en.wikipedia.org/wiki/Reproducible_builds)

If someone is to download the packages for this project by running **go mod download** they will get the exact versions of all the packages the project and can verify the integrity of the packages:

```sh
    go mod verify
```

The above command will verify that the checksums of the downloaded packages on your machine match the entries in **go.sum**, so you can be confident that they haven’t been altered.

## Templating 

The default **html/template** Go library was used as a tempting engine. The default library automatically escapes any data yielded between {{ }}, it is also smart enough to make escaping context-dependent.

It will use the appropriate escape sequences depending on whether the data is rendered in a part of the page that contains HTML, CSS, Javascript or a URI.


## Middleware 

Middlewares are useful when you want to share some functionality across multiple HTTP requests e.g. you might want to log every request, compress every response, or check a cache before passing the request to your handlers.

A middleware essentially is a self-contained code which independently acts on a request before or after your normal application handlers.

Think of a Go web application as a chain of **ServeHTTP()** methods being called
one after another e.g. HTTPRequest -> ServerMux's ServerHTTP() -> relevant handler's ServeHTTP(). The basic idea of middleware is to insert another handler into this chain think **http.StripPrefix()** function from serving static files which removes a specific prefix from the request’s URL path before passing the request on to the file server. 

This project includes the following middleware: 

    Before the ServeMux:
    1. secureHeaders - adds CSP headers
    2. logRequest - logs the details of each incoming request
    3. recoverPanic - takes control of a panic and handlers it gracefully for the user

Go ServeHttp handles well panics and doesn't bring down the web application, we also handle it with receiverPanic middleware but you would have to handle independently each goroutine you spin up e.g.

```go
    func (app *application) myHandler(w http.ResponseWriter, r *http.Request) {
    
        // Spin up a new goroutine to do some background processing.
        go func() {
            defer func() {
                if err := recover(); err != nil {
                    app.logger.Error(fmt.Sprint(err))
                }
            }()

            doSomeBackgroundProcessing()
        }()

        w.Write([]byte("OK"))
    }
```

## Routing 

Go’s servemux doesn’t support method based routing or clean URLs with variables in them which is why it's more appropriate to chose a third-party package than to build a custom solution. 

The most [popular and well supported routers](https://www.alexedwards.net/blog/which-go-router-should-i-use) are: 
1. **julienschmidt/httprouter** - the most focused, lightweight and fastest of the three packages. It automatically handles OPTIONS requests and sends 405 responses correctly, and allows you to set custom handlers for 404 and 405 responses too.

2. **go-chi/chi** - similar to httprouter in terms of its features, with the main differences being that it also supports regexp route patterns and ‘grouping’ of routes which use specific middleware. The only real
downside of chi is that it doesn’t automatically handle OPTIONS requests for you.

3. **gorilla/mux** - is the most full-featured router it supports regexp route patterns, and allows you to route requests based on scheme, host and headers. It’s also the only one to support custom routing rules and route ‘reversing’. The downsides of ***gorilla/mux*** is that it’s comparatively slow and memory hungry and like ***chi***, it also doesn’t automatically handle OPTIONS requests, plus it doesn’t set an Allow header in 405 responses

Thus, for the sake of performance and correctness, for this project ***julienschmidt/httprouter*** was used. 

```sh
     go get github.com/julienschmidt/httprouter@v1
```

