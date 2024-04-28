# orthography

**orthography** - this is a library for working with Russian orthography

Install:
```shell
go get github.com/artarts36/orthography
```

## Incline

### Incline directly through morpher.ru

```go
package main

import (
	"context"
	"fmt"
	"github.com/artarts36/orthography/incline"
)

func main() {
	words, _ := incline.NewMorpherDefault().InclineNouns(context.Background(), []string{
		"Стол",
		"Стул",
	})
	for _, word := range words {
		fmt.Println(word.Nominative, word.Genitive, word.Dative, word.Accusative, word.Instrumental, word.Prepositional)
	}
}
```

### Incline directly through morpher.ru

```go
package main

import (
	"context"
	"fmt"
	"github.com/artarts36/orthography/incline"
)

func main() {
	words, _ := incline.NewMorpherDefault().InclineNouns(context.Background(), []string{
		"Стол",
		"Стул",
	})
	for _, word := range words {
		fmt.Println(word.Nominative, word.Genitive, word.Dative, word.Accusative, word.Instrumental, word.Prepositional)
	}
}
```

### Incline with caching in-memory

```go
package main

import (
	"context"
	"fmt"
	"github.com/artarts36/orthography/incline"
	"github.com/artarts36/orthography/incline/store"
)

func main() {
	inc := incline.NewPersistent(
		incline.NewMorpherDefault(),
		store.NewMemory(),
	)

	res, _ := inc.InclineNouns(context.Background(), []string{"Стол", "Стул"})
	for _, word := range words {
		fmt.Println(word.Nominative, word.Genitive, word.Dative, word.Accusative, word.Instrumental, word.Prepositional)
	}
}
```

### Incline with persist in database

```go
package main

import (
	"context"
	"fmt"
	"github.com/artarts36/orthography/incline"
	"github.com/artarts36/orthography/incline/store"
)

func main() {
	db, _ := sql.Open("postgres", "host=localhost port=5222 user=test password=test dbname=orthography sslmode=disable")
	
	inc := incline.NewPersistent(
		incline.NewMorpherDefault(),
		store.NewDB("words", db),
	)

	res, _ := inc.InclineNouns(context.Background(), []string{"Стол", "Стул"})
	for _, word := range words {
		fmt.Println(word.Nominative, word.Genitive, word.Dative, word.Accusative, word.Instrumental, word.Prepositional)
	}
}
```


### Incline with caching in-memory and persist in database

```go
package main

import (
	"context"
	"fmt"
	"github.com/artarts36/orthography/incline"
	"github.com/artarts36/orthography/incline/store"
)

func main() {
	db, _ := sql.Open("postgres", "host=localhost port=5222 user=test password=test dbname=orthography sslmode=disable")
	
	inc := incline.NewPersistent(
		incline.NewMorpherDefault(),
		store.NewProxyMemoryAndDB("words", db),
	)

	res, _ := inc.InclineNouns(context.Background(), []string{"Стол", "Стул"})
	for _, word := range words {
		fmt.Println(word.Nominative, word.Genitive, word.Dative, word.Accusative, word.Instrumental, word.Prepositional)
	}
}
```

Also, warmup memory store

```go
package main

import (
	"context"
	"fmt"
	"github.com/artarts36/orthography/incline"
	"github.com/artarts36/orthography/incline/store"
)

func main() {
	db, _ := sql.Open("postgres", "host=localhost port=5222 user=test password=test dbname=orthography sslmode=disable")
	
	st := store.NewProxyMemoryAndDB("words", db)
	
	inc := incline.NewPersistent(
		incline.NewMorpherDefault(),
		st,
	)
	
	st.WarmUp()

	res, _ := inc.InclineNouns(context.Background(), []string{"Стол", "Стул"})
	for _, word := range words {
		fmt.Println(word.Nominative, word.Genitive, word.Dative, word.Accusative, word.Instrumental, word.Prepositional)
	}
}
```
