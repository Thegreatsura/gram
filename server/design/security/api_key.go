package security

import (
	"fmt"

	"github.com/speakeasy-api/gram/server/internal/auth"
	. "goa.design/goa/v3/dsl"
)

var ByKey = APIKeySecurity(auth.KeySecurityScheme, func() {
	Description("key based auth.")
	Scope("consumer", "consumer based tool access")
	Scope("producer", "producer based tool access")
})

var ByKeyPayload = func() {
	APIKey(auth.KeySecurityScheme, "apikey_token", String)
}

var ByKeyHeader = func() {
	Header(fmt.Sprintf("apikey_token:%s", auth.APIKeyHeader), String, "API Key header")
}

var ByKeyNamedHeader = func(name string) {
	Header(fmt.Sprintf("apikey_token:%s", name), String, "API Key header")
}
