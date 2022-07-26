# DI - Address Issuer Repo

Digital Identity is the identity use case of verifiable credentials as a implementation of the technology of HyperLedger Indy, HyperLedger URSA and Hyperledger Aries.

- This use case involves issuing a physical address digital credential.

## Directory structure

```
.
├── api             // handlers, router (endpoints) and ui (production build)
├── cmd             // server entry point
├── deploy          // k8s aws deployment files
├── cmd             // entry point
├── pkg
│   ├── client      // ACA-py client functions
│   ├── config      // config loading env variables
│   ├── log         // logger
│   ├── models      // all dto's
│   ├── server      // server configuration and middleware
│   └── utils       // helper functions
├── web             // React.js ui
└─
```
