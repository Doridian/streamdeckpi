module github.com/Doridian/streamdeckpi

go 1.19

require github.com/Doridian/streamdeck v0.1.2

require (
	github.com/KarpelesLab/hid v0.1.0 // indirect
	golang.org/x/image v0.0.0-20220902085622-e7cb96979f69 // indirect
)

// Well this seems like a weird kludge, but it works to make gokrazy understand to use the local files...
replace github.com/Doridian/streamdeckpi => ./
