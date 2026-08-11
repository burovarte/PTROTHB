package main

import (
	"fmt"

	"github/burovarte/PTROTHB/sprint_2/cli"
	patterns "github/burovarte/PTROTHB/sprint_3/patterns"
)

func main() {
	fmt.Println("Hello, World!")

	patterns.Pipeline()
	patterns.MainTimeout()
	patterns.MainOrDone()
	patterns.MainErrGroup()
	patterns.MainRateLimiting()
	patterns.MainOrChanel()
	patterns.MainBridge()

	cli.Run()
}
