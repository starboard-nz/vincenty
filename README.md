**vincenty**

Calculate the distance between 2 points on Earth's surface using Vincenty's formula.

This is a Go re-implementation of the vincenty Python library.

**usage**

```
package main

import (
        "fmt"
        "github.com/csst/common/vincenty"
)

func main() {
        Dunedin := vincenty.LatLng{Latitude: -45.8726082, Longitude: 170.3870355 }
        Alexandra := vincenty.LatLng{Latitude: -45.2426447, Longitude: 169.0977066 }

	dist := vincenty.Inverse(Dunedin, Alexandra)
	fmt.Printf("Distance: %vkm\n", dist.Kilometers())
}
```

**References**

* https://github.com/maurycyp/vincenty/
* https://en.wikipedia.org/wiki/Vincenty's_formulae
* https://en.wikipedia.org/wiki/World_Geodetic_System
