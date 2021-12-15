module github.com/KlyuchnikovV/buffer

go 1.16

require (
	github.com/KlyuchnikovV/edigode v0.0.0-20211209121658-42980a9c6cd7
	github.com/KlyuchnikovV/gapbuf v0.0.0-20211209120053-9298a9e329ad
	github.com/KlyuchnikovV/linetree v0.0.0-20211209120521-7dcdde3ffe70
	github.com/wailsapp/wails v1.16.8
)

replace (
	github.com/KlyuchnikovV/edigode => ../edigode
	github.com/KlyuchnikovV/gapbuf => ../gapbuf
	github.com/KlyuchnikovV/linetree => ../linetree
)
