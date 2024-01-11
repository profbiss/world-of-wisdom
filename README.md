# "Word of Wisdom" tcp server (protected from DDOS attacks with the Proof of Work).

## Task

Design and implement "Word of Wisdom" tcp server:

- TCP server should be protected from DDOS attacks with the Prof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
- The choice of the POW algorithm should be explained.
- After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
- Docker file should be provided both for the server and for the client that solves the POW challenge.

## Getting started

Requirements:
- Docker installed (to run docker-compose)

```
# Run server and client by docker-compose
cd build && docker-compose up
```

## Resources

- [Word of Wisdom](https://en.wikipedia.org/wiki/Word_of_Wisdom)
- [Proof of work](https://en.wikipedia.org/wiki/Proof_of_work)
- [Hashcash](https://en.wikipedia.org/wiki/Hashcash)
- [Distributed Consensus – Proof-of-Work](https://oliverjumpertz.com/distributed-consensus-proof-of-work/)