module github.com/common-fate/ciem

go 1.21.3

require (
	github.com/AlecAivazis/survey/v2 v2.3.7
	github.com/briandowns/spinner v1.23.0
	github.com/bufbuild/connect-go v1.10.0
	github.com/charmbracelet/bubbles v0.17.1
	github.com/charmbracelet/bubbletea v0.25.0
	github.com/charmbracelet/huh v0.0.0
	github.com/charmbracelet/huh/spinner v0.0.0
	github.com/charmbracelet/lipgloss v0.9.1
	github.com/common-fate/clio v1.2.3
	github.com/common-fate/grab v1.2.0
	github.com/common-fate/sdk v0.0.0-20231220055240-32a937c74bf1
	github.com/fatih/color v1.16.0
	github.com/mattn/go-isatty v0.0.20
	github.com/mattn/go-runewidth v0.0.15
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.8.4
	github.com/urfave/cli/v2 v2.25.7
	go.uber.org/zap v1.26.0
	golang.org/x/oauth2 v0.14.0
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/alecthomas/chroma v0.10.0 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/catppuccin/go v0.2.0 // indirect
	github.com/charmbracelet/glamour v0.6.0 // indirect
	github.com/dlclark/regexp2 v1.10.0 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/microcosm-cc/bluemonday v1.0.25 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/yuin/goldmark v1.6.0 // indirect
	github.com/yuin/goldmark-emoji v1.0.2 // indirect
)

require (
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.2 // indirect
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/atotto/clipboard v0.1.4 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/charmbracelet/x/exp/strings v0.0.0-20231215171016-7ba2b450712d
	github.com/common-fate/apikit v0.3.0 // indirect
	github.com/containerd/console v1.0.4-0.20230313162750-1ae8d489ac81 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/danieljoos/wincred v1.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dvsekhvalnov/jose2go v1.5.0 // indirect
	github.com/fatih/structtag v1.2.0 // indirect
	github.com/go-chi/chi/v5 v5.0.10 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.4.0 // indirect
	github.com/gorilla/schema v1.2.0 // indirect
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/lithammer/fuzzysearch v1.1.8
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/muesli/termenv v0.15.2 // indirect
	github.com/muhlemmer/gu v0.3.1 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	github.com/zitadel/oidc/v2 v2.12.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/crypto v0.15.0 // indirect
	golang.org/x/net v0.18.0 // indirect
	golang.org/x/sync v0.5.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/term v0.14.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/charmbracelet/huh v0.0.0 => ../huh

replace github.com/charmbracelet/accessibility v0.0.0 => ../huh/accessibility

replace github.com/charmbracelet/huh/spinner v0.0.0 => ../huh/spinner
