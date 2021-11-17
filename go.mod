module github.com/KlyuchnikovV/buffer

go 1.16

require (
	github.com/KlyuchnikovV/gapbuf v0.0.0-20211103104016-70b130518d59
	github.com/KlyuchnikovV/linetree v0.0.0-20211106122630-5decd8ddd752
	github.com/stretchr/testify v1.7.0
)

replace (
	github.com/KlyuchnikovV/gapbuf => ../gapbuf
	github.com/KlyuchnikovV/linetree => ../linetree
)
