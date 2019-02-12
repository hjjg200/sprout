package sprout

type Authenticator func( *Request ) bool