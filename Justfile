set shell := ["nu", "-c"]
set dotenv-load

default:
	@just --choose

# using https://github.com/cosmtrek/air - global install
dev:
	air

run:
	go run .

# install via templ
templ:
  templ generate -watch $"-proxy=http://localhost:($env.PORT)"

# install via pnpm i -g tailwindcss
css:
	tailwindcss -i static/style.css -o static/style-compiled.css --minify --watch

# generate db from schema changes
ent:
	go generate ./src/database/ent
