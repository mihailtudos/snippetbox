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

A better but more complex approach is to create a custom implementation of
http.FileSystem, and have it return an os.ErrNotExist. 