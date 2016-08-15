/*Package transforms imports all of the transforms that are available with PipeScript. The core PipeScript
only has an if statement and the identity operator, which are not nearly enough.

This package imports EVERYTHING
*/
package transforms

import (
	"github.com/connectordb/pipescript/transforms/core"     // The core transforms
	"github.com/connectordb/pipescript/transforms/datetime" // Manipulating timestamps
	"github.com/connectordb/pipescript/transforms/math"     // Statistical transforms
	"github.com/connectordb/pipescript/transforms/misc"     // Miscellaneous transforms
	"github.com/connectordb/pipescript/transforms/strings"  // Text-based transforms
)

// Register ALL functions
func Register() {
	core.Register()
	math.Register()
	datetime.Register()
	strings.Register()
	misc.Register()
}
