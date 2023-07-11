# go-encoding-form

Marshaling and Unmarshaling x-www-form-urlencoded contents in Go.

> **Warning**  
> This is not an official product of YUMEMI Inc.


## Usages

### Marshaling Form Values

```go
package main

import (
	"fmt"
	
	"github.com/yumemi-inc/go-encoding-form"
)

type FormValues struct {
	Foo string `form:"foo"`
	Bar *string `form:"bar,omitempty"`
}

func main() {
	v := &FormValues{
		Foo: "Lorem ipsum",
		Bar: nil,
	}

	bytes, _ := form.MarshalForm(v)
	fmt.Printf("%s\n", string(bytes))
	// foo=Lorem+ipsum
}
```

### Unmarshaling Form Values

```go
package main

import (
	"fmt"

	"github.com/yumemi-inc/go-encoding-form"
)

type FormValues struct {
	Foo string  `form:"foo"`
	Bar *string `form:"bar,omitempty"`
}

func main() {
	v := new(FormValues)
	_ = form.UnmarshalForm([]byte("foo=Lorem+ipsum"), v)
	fmt.Printf("%+v\n", v)
	// &{Foo:Lorem ipsum Bar:<nil>}
}
```
