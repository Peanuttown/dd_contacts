module github.com/Peanuttown/dd_contacts

go 1.15

require (
	entgo.io/ent v0.8.0
	github.com/Peanuttown/dd_api v0.0.9
	github.com/Peanuttown/tzzGoUtil v1.2.2
	github.com/go-sql-driver/mysql v1.5.1-0.20200311113236-681ffa848bae
)

replace github.com/Peanuttown/tzzGoUtil v1.2.2 => ../tzzGoUtil

replace github.com/Peanuttown/dd_api v0.0.9 => ../dd_api
