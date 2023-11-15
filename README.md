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

To serve static file the http.FileServe handler was used. It has multiple advantages:

1. It sanitizes all requests to stop directory traversal attacks
2. Range requests are fully supported (support resumable downloads)
3. The Last-Modified and If-Modified-Since headers are transparently supported
4. The Content-Type is automatically set from the file extension u

To serve single file a new handler can be set as (warning: http.ServeFile() does not automatically sanitize the file path - you must sanitize the input with filepath.Clean() before using it): 

```go
func downloadHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./ui/static/file.zip")
}
```

Disable static file by creating index.html files in each static subdirectory:

```sh
    find ./ui/static -type d -exec touch {}/index.html \;
```

## Logging - structured logger

To decouple the logging from the application a basic structured, concurrency-safe, logger was created with the default slog package (set to: os.Stdout). The final destination of the logs can be managed by your execution environment independently of the application: 

```sh 
    go run ./cmd/web >> /tmp/web.log
```

## DI 

This web application, as most web applications, has multiple dependencies that the handlers need to access, such as a database connection pool, centralized error handlers, and template caches. Thus, to make these dependencies available, following the good practices they will be inject. 

DI is build around the application struct and works well with the handlers as long as they are in the same package **main**, for a more complex project structure [Closures for dependency injection](https://gist.github.com/alexedwards/5cd712192b4831058b21) can be used. 

A better but more complex approach is to create a custom implementation of
http.FileSystem, and have it return an os.ErrNotExist. 