# unearth

Unearth unadvertised job positions by crawling GitHub organizations' web pages!

## Architecture

`pull.go` -> Pulls organizations using GitHub API and saves them in an SQLite database.

`descend.go` -> Follows "Careers" (or similar) links from the main page to positions listing.

`unearth.go` -> Extracts positions and locations from the listings page using OCRand saves in the database.

## Run

`go build && ./unearth [pull] [process]`

Application accepts two arguments of string literals `pull` and `process`.

Passing in `pull` without `process` will just fetch the organizations and save them in the database for later processing.

Passing in `process` without `pull` will descend/unearth the organizations that are in the database.

Passing in both will fetch organizations and descend/unearth them at the same time.
