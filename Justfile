set shell := ["nu", "-c"]

default:
	@just --choose

# using https://github.com/cosmtrek/air - global install
dev:
	air

run:
	go run .

css: 
	pnpx tailwindcss -i static/style.css -o static/style-compiled.css --minify
