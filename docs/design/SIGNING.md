# Signing

A core aspect of the system is message signing. Using [public key cryptography](https://en.wikipedia.org/wiki/Public-key_cryptography),
clients should be able to sign the messages they produce so that consumers can
verify the messages were sent by that client.

To achieve this, servers can generate a key pair on demand for clients which
are immutable. Once a client has obtained a key pair it _must_ persist the
private key somewhere, as the server will not store it. Instead, the server
will store the public key and use that to verify the signed message from a
client when it is produced.

Servers use [Ed25519](https://en.wikipedia.org/wiki/EdDSA) to generate keys,
clients _must_ use the same for signing their messages. Specifically, the
[nacl/sign](https://pkg.go.dev/golang.org/x/crypto/nacl/sign) package is used
for go based client/server implementations

For efficiency, it is not advised to sign large messages. To avoid this, every
message has the concept of a key and value. The message value should contain
the content of the message to be read by a consumer. The message key should
contain some metadata regarding the message value, that can be encoded and
signed and provided alongside the message itself to the server. Without
providing a message key, the message itself cannot be verified by the server.

When creating a new topic, it is possible to enable "verified-only"
messages. This means that the server will reject any attempts to produce a
message onto that topic without a verifiable message key. This allows operators
to decide if they require message signing on particular topics.

The implementation for this can be seen in the `internal/signing` package
within the source code.
