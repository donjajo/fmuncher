# fmuncher
Takes file pointer, munch it! 
The idea is to manipulate files as much as possible. Currently, the only feature is to split files into parts of offset and length.

## Installation
Using `go get`
```sh
    $ go get github.com/donjajo/fmuncher
```
## Usage
    ```
        file, fileErr := os.Open(filename)
        if fileErr != nil {
            log.Fatal(fileErr)
        }

        muncher := fmuncher.Munch(file)
        fmt.Println(muncher.Split())
    ```

More doc on [Godoc](https://godoc.org/github.com/donjajo/fmuncher)

## Contribution
    * Have a experience with better documentation. Well accepted!
    * Have an idea? Please raise an issue :)