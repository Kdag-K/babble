// Package keys implements the public key cryptography used throughout Kdag.
//
// An instance of a Kdag node, also referred to as peer, participant or
// validator, owns a cryptographic key-pair that it uses to encrypt, sign and
// verify messages. The private key is secret but the public key is used by
// other nodes to verify messages signed with the private key.
//
// Kdag uses elliptic curve cryptography (ECDSA) with the sec256k1 curve. We
// chose the secp256k1 curve because it is also used by Bitcoin and Ethereum
// which means that Bitcoin and Ethereum keys can be used to operate a Kdag
// node.r
package keys
