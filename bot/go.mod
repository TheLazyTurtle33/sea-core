module github.com/TheLazyTurtle33/sea-core/bot

go 1.26.1

require github.com/TheLazyTurtle33/sea-core/shared v0.0.0

replace github.com/TheLazyTurtle33/sea-core/shared => ../shared

require (
	github.com/lib/pq v1.12.0 // indirect
	golang.org/x/net v0.52.0 // indirect
)
