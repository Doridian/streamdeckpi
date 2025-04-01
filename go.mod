module github.com/Doridian/streamdeckpi

go 1.23.0

toolchain go1.24.2

require (
	github.com/Doridian/go-haws v0.4.3
	github.com/Doridian/go-streamdeck v1.4.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	golang.org/x/image v0.25.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/Doridian/karalabe_hid v1.0.0 // indirect
	github.com/KarpelesLab/hid v0.1.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
)

// Well this seems like a weird kludge, but it works to make gokrazy understand to use the local files...
replace github.com/Doridian/streamdeckpi => ./

// replace github.com/Doridian/go-haws => ../go-haws/
// replace github.com/Doridian/go-streamdeck => ../go-streamdeck/
