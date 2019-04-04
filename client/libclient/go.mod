module libclient

go 1.12

replace golang.org/x/exp => github.com/golang/exp v0.0.0-20190321205749-f0864edee7f3

replace golang.org/x/tools => github.com/golang/tools v0.0.0-20190325223049-1d95b17f1b04

replace golang.org/x/net => github.com/golang/net v0.0.0-20190324223953-e3b2ff56ed87

replace golang.org/x/sys => github.com/golang/sys v0.0.0-20190322080309-f49334f85ddc

replace golang.org/x/mobile => github.com/golang/mobile v0.0.0-20190319155245-9487ef54b94a

replace golang.org/x/image => github.com/golang/image v0.0.0-20190321063152-3fc05d484e9f

replace golang.org/x/text => github.com/golang/text v0.3.0

replace golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190325154230-a5d413f7728c

require (
	github.com/hajimehoshi/oto v0.3.1
	github.com/rocket049/gettext-go v0.0.0-20190404080233-af421a50b332
	github.com/russross/blackfriday v1.5.2
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/skratchdot/open-golang v0.0.0-20190104022628-a2dfa6d0dab6
)
