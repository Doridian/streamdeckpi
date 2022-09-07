module github.com/Doridian/streamdeckpi

go 1.19

require github.com/muesli/streamdeck v0.3.0

require (
	github.com/karalabe/hid v1.0.1-0.20190806082151-9c14560f9ee8 // indirect
	golang.org/x/image v0.0.0-20200119044424-58c23975cae1 // indirect
)

// Well this seems like a weird kludge, but it works to make gokrazy understand to use the local files...
replace github.com/Doridian/streamdeckpi => ./
