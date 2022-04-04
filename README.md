# go-wfc-algorithm

Animated example showing randomly generated constraint based tile maps using the `go-wfc` package.

View the browser build here:
https://zfedoran.github.io/go-wfc-algorithm/


## Overview
This package uses the *Wave Function Collapse* algorithm as described by Oskar
Stålberg.

The wave function collapse algorithm is a recursive algorithm that picks a
random tile for a slot on the output image and removes impossible neighbors
until only a single possibility remains. The algorithm and more details are 
described on the [go-wfc](https://github.com/zfedoran/go-wfc) repository.

## Build

The project uses [ebiten](https://ebiten.org/) have a look at the WASM build
steps here: https://ebiten.org/documents/webassembly.html


## Artwork

The awesome artwork in this repository was done by
[@makionfire](https://twitter.com/makionfire). If you need help designing a tile
set, I highly recommend reaching out to her. A huge shout-out to `@makionfire`
for letting me use this tileset.

The artwork itself does **not** fall under the MIT licence.

## Licence

The licence for the source code in this package is MIT. Meaning, do whatever
you'd like but we'd love a shoutout. The goal is to get more folks to build
games with golang.

-------------------------------------------------------------------------------

If you like this work and want to buy me or the artist a coffee or beer, you're
free to do so by sending to some SOL to
[🍺💵.sol](https://naming.bonfida.org/#/domain/%F0%9F%8D%BA%F0%9F%92%B5)
