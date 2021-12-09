module github.com/KlyuchnikovV/buffer

go 1.16

require (
	github.com/KlyuchnikovV/edigode v0.0.0-20211209121658-42980a9c6cd7
	github.com/KlyuchnikovV/gapbuf v0.0.0-20211209120053-9298a9e329ad
	github.com/KlyuchnikovV/linetree v0.0.0-20211209120521-7dcdde3ffe70
	github.com/davecgh/go-spew v1.1.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace (
	github.com/KlyuchnikovV/gapbuf => ../gapbuf
	github.com/KlyuchnikovV/linetree => ../linetree
)
