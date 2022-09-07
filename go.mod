module github.com/Doridian/streamdeckpi

go 1.19

require github.com/Doridian/streamdeck v0.0.0-20220907015624-9bf18f1d511d

require (
	github.com/jkassis/hid v0.0.0-20220630003547-398145ff2de0 // indirect
	golang.org/x/image v0.0.0-20200119044424-58c23975cae1 // indirect
)

// Well this seems like a weird kludge, but it works to make gokrazy understand to use the local files...
replace github.com/Doridian/streamdeckpi => ./
