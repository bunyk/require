The purpose of this module is to allow you to "require" text files so they could be hardcoded in single executable and all filesystem access happens during the build. 

# Usage

```go
import "github.com/bunyk/require"

...
text := require.Require("file.txt")
...
```

