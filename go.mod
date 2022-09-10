module github.com/Doridian/streamdeckpi

go 1.19

require (
	github.com/Doridian/go-haws v0.3.1
	github.com/Doridian/go-streamdeck v1.1.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/KarpelesLab/hid v0.1.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/karalabe/hid v1.0.1-0.20190806082151-9c14560f9ee8 // indirect
	golang.org/x/image v0.0.0-20220902085622-e7cb96979f69 // indirect
)

// Well this seems like a weird kludge, but it works to make gokrazy understand to use the local files...
replace github.com/Doridian/streamdeckpi => ./

// replace github.com/Doridian/go-haws => ../go-haws/
// replace github.com/Doridian/go-streamdeck => ../streamdeck/
