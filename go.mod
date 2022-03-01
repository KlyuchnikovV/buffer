module github.com/KlyuchnikovV/buffer

go 1.16

require (
	github.com/KlyuchnikovV/edigode v0.0.0-20211209121658-42980a9c6cd7
	github.com/KlyuchnikovV/gapbuf v0.0.0-20211209200800-ff61a486347e
	github.com/KlyuchnikovV/linetree v0.0.0-20211209200702-afa8ac3d48ba
	golang.design/x/clipboard v0.6.0
)

replace (
	github.com/KlyuchnikovV/edigode => ../edigode
	github.com/KlyuchnikovV/gapbuf => ../gapbuf
	github.com/KlyuchnikovV/linetree => ../linetree
)
