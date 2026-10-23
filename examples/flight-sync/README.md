# Example Flight Sync Engine

Example engine whose processes share one credential per Terragrunt invocation. Terragrunt starts one engine process per unit, so without coordination a `terragrunt run --all` across ten units exchanges a credential ten times.

Terragrunt calls one invocation a flight, and sends two fields in every `InitRequest` of the flight:

- `flight_id`: a UUIDv7, the same for every engine process in the flight. It is not secret.
- `flight_token`: 32 random bytes, the same for every engine process in the flight. It is secret.

The engine uses them like this:

1. It opens a directory named after `flight_id` under its state directory. Every engine in the flight opens the same one.
2. It takes an exclusive lock on `credential.lock` in that directory. The other engines wait.
3. It reads `credential`, which holds the credential encrypted with AES-GCM under a key derived from `flight_token`. If the file is missing, expired, or fails to decrypt, the engine exchanges a new credential with the auth server and writes it back encrypted.
4. It releases the lock, and the next engine finds the stored credential.

The first engine to take the lock exchanges, and the rest reuse its credential. If Terragrunt sends no flight, because it predates these fields, each engine exchanges a credential for itself.

The [flight](./flight) package holds the lock and encryption code.

Nothing removes old flight directories. A real engine would remove the ones whose UUIDv7 timestamp is old.

## Building

```bash
make
```

This builds `terragrunt-engine-flight-sync`, the engine, and `flight-sync-auth`, a stand-in auth server that returns a random credential for every exchange and prints a running count.

## Running

1. Start the auth server:

```bash
./flight-sync-auth
```

2. In another shell, run every unit in `live`:

```bash
cd live
export TG_EXPERIMENTAL_ENGINE=1
export TG_ENGINE_SKIP_CHECK=1
terragrunt run --all plan
```

Each unit logs which credential it used. With a Terragrunt that sends flights, one unit exchanges and the others reuse its credential, and the auth server prints only `exchange 1`:

```text
[a] Exchanged credential 1f2e3d4c and stored it for the flight
[b] Reused the flight's credential 1f2e3d4c
[c] Reused the flight's credential 1f2e3d4c
```

A Terragrunt release from before flights sends neither field, so every unit exchanges its own credential, and the auth server prints one exchange per unit:

```text
[a] Exchanged credential 9bcc99e6 for this engine alone, because Terragrunt sent no flight
[b] Exchanged credential c22e405f for this engine alone, because Terragrunt sent no flight
[c] Exchanged credential 7d21a783 for this engine alone, because Terragrunt sent no flight
```

`engine/main_test.go` starts engine processes the way Terragrunt does and sends them a shared flight, so `make test` covers the first case.

## Configuration

The engine reads these from the environment Terragrunt sends in `InitRequest`:

- `FLIGHT_SYNC_AUTH_URL`: where to exchange credentials. Defaults to `http://127.0.0.1:8080/exchange`.
- `FLIGHT_SYNC_STATE_DIR`: where to create flight directories. Defaults to `terragrunt-engine-flight-sync` under the user's cache directory.
