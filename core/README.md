# Zeta core

A decentralised trading platform that allows pseudo-anonymous trading of derivatives on a blockchain.

**Zeta** provides the following core features:

- Join a Zeta network as a validator or non-consensus node.
- [Governance](./governance/README.md) - proposing and voting for new markets
- A [matching engine](./matching/README.md)
- [Configure a node](#configuration) (and its [APIs](#apis))
- Manage authentication with a network
- [Run scenario tests](./integration/README.md)
- Support settlement in cryptocurrency (coming soon)
## Links

- For **new developers**, see [Getting Started](../GETTING_STARTED.md).
- For **updates**, see the [Change log](../CHANGELOG.md) for major updates.
- For **architecture**, please read the [documentation](./docs#zeta-core-architecture) to learn about the design for the system and its architecture.
- Please [open an issue](https://github.com/zetaprotocol/zeta/issues/new) if anything is missing or unclear in this documentation.

<details>
  <summary><strong>Table of Contents</strong> (click to expand)</summary>

<!-- toc -->

- [Zeta](#zeta-core)
  - [Links](#links)
  - [Installation](#installation)
  - [Configuration](#configuration)
    - [Files location](#files-location)
  - [Zeta node wallets](#zeta-node-wallets)
    - [Using Ethereum Clef wallet](#using-ethereum-clef-wallet)
      - [Automatic approvals](#automatic-approvals)
      - [Importing and generation account](#importing-and-generation-account)
  - [API](#api)
  - [Provisioning](#provisioning)
  - [Troubleshooting & debugging](#troubleshooting--debugging)

<!-- tocstop -->

</details>

## Installation

To install `trading-core` and `tendermint`, see [Getting Started](../GETTING_STARTED.md).

## Configuration

Zeta is initialised with a set of default configuration with the command `zeta init`. To override any of the defaults, edit your `config.toml`.

**Example**

```toml
[Matching]
Level = 0
ProRataMode = false
LogPriceLevelsDebug = false
LogRemovedOrdersDebug = false
```

Zeta requires a set of wallets for the internal or external chain it's dealing with.

The node wallets can be accessed using the `nodewallet` subcommand, these node wallets are initialized / accessed using a passphrase that needs to be specified when initializing Zeta:

```sh
zeta init --nodewallet-passphrase-file "my-passphrase-file.txt"
```

### Files location

| Environment variables | Unix             | MacOS                           | Windows                |
| :-------------------- | :----------------| :------------------------------ | :--------------------- |
| `XDG_DATA_HOME`       | `~/.local/share` | `~/Library/Application Support` | `%LOCALAPPDATA%`       |
| `XDG_CONFIG_HOME`     | `~/.config`      | `~/Library/Application Support` | `%LOCALAPPDATA%`       |
| `XDG_STATE_HOME`      | `~/.local/state` | `~/Library/Application Support` | `%LOCALAPPDATA%`       |
| `XDG_CACHE_HOME`      | `~/.cache`       | `~/Library/Caches`              | `%LOCALAPPDATA%\cache` |

You can override these environment variables, however, bear in mind it will apply system-wide.

If you don't want to rely on the default XDG paths, you can use the `--home` flag on the command-line.

## Zeta node wallets

A Zeta node needs to connect to other blockchain for various operation:
- validate transaction happened on foreign chains
- verify presence of assets
- sign transaction to be verified on foreign blockchain
- and more...

In order to do these different action, the Zeta node needs to access these chains using their native wallet. To do so the zeta command line provide a command line tool:
`zeta nodewallet` allowing users to import foreign blockchain wallets credentials, so they can be used at runtime.

For more details on how to use the Zeta node wallets run:
```
zeta nodewallet --help
```

### Using Ethereum Clef wallet

#### Automatic approvals

Given that Clef requires manually approving all RPC API calls, it is mandatory to setup
[custom rules](https://github.com/ethereum/go-ethereum/blob/master/cmd/clef/rules.md#rules) for automatic approvals. Zeta requires at least `ApproveListing` and `ApproveSignData` rules to be automatically approved.

Example of simple rule set JavaScript file with approvals required by Zeta:
```js
function ApproveListing() {
  return "Approve"
}

function ApproveSignData() {
  return "Approve"
}
```

Clef also allows more refined rules for signing. For example approves signs from `ipc` socket:
```js
function ApproveSignData(req) {
  if (req.metadata.scheme == "ipc") {
    return "Approve"
  }
}
```

Please refer to Clef [rules docs](https://github.com/ethereum/go-ethereum/blob/master/cmd/clef/rules.md#rules) for more information.

#### Importing and generation account

As of today, Clef does not allow to generate a new account for other back end storages than a local Key Store. Therefore it is preferable to create a new account on the back end of choice and import it to Zeta through node wallet CLI.

Example of import:
```sh
zeta nodewallet import --chain=ethereum --eth.clef-address=http://clef-address:port
```


## Troubleshooting & debugging

The application has structured logging capability, the first port of call for a crash is probably the Zeta and Tendermint logs which are available on the console if running locally or by journal plus syslog if running on test networks. Default location for log files:

* `/var/log/zeta.log`
* `/var/log/tendermint.log`

Each internal Go package has a logging level that can be set at runtime by configuration. Setting the logging `Level` to `"Debug"` for a package will enable all debugging messages for the package which can be useful when trying to analyse a crash or issue.

Debugging the application locally is also possible with [Delve](../DEBUG_WITH_DLV.md).
