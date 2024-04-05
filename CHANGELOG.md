# @common-fate/cli

## 1.12.1

### Patch Changes

- 2c79b0d: Fixes a regression in the Keychain which meant switching contexts would require re-authenticating

## 1.12.0

### Minor Changes

- 922c95f: Adds support for specifying config sources via the CF_CONFIG_SOURCES environment variable. You can now use `CF_CONFIG_SOURCES=env` in CI environments to configure the CLI entirely via environment variables.

## 1.11.0

### Minor Changes

- c1f0a8a: Allows environment variables to be used to configure the CLI.

  - `CF_OIDC_ISSUER`: specifies the OIDC issuer

  - `CF_OIDC_CLIENT_ID`: specifies the OIDC client ID

  - `CF_OIDC_CLIENT_SECRET`: specifies the OIDC client secret

  - `CF_API_URL` specifies the API URL

  - `CF_ACCESS_URL` specifies the Access URL (the URL that the Common Fate Access service is hosted on)

  - `CF_AUTHZ_URL` specifies the Authz URL (the URL that the Common Fate Authz service is hosted on)

- f28e678: Adds new commands for Access Simulation APIs (which are currently in beta).

  cf access list approvers
  cf access preview user-access
  cf access preview entitlement-access

## 1.10.0

### Minor Changes

- 7b1e038: Shifts 'cf policyset' commands to be 'cf authz policyset'.
- 7b1e038: Adds 'cf authz policyset validate' command to validate Cedar policies

## 1.9.0

### Minor Changes

- dc6b5e3: Adds entra-relink-users command which can be used to remediate issues with entra scim configurations which cause users not to be linked correctly
- 85ede0e: added optional 'duration' flag to allow override of default access duration when using ensure request
- 37b1784: Adds 'cf authz schema get' command which will download the Cedar schema used for authorization.
- 7747009: list availabilities will now display deduplicated target roles by default

### Patch Changes

- 2dd3fc6: Adds cli command for retrying background jobs if they are retryable. Also adds timestamp to background job list command
- 1b16429: Fixes the output formatting for the 'cf deployment logs get' command.

## 1.8.0

### Minor Changes

- dd7da7d: Add command to reset entra users in the database so that they may be re-synced

## 1.7.0

### Minor Changes

- 0513f58: Adds background task diagnostic command for listing background tasks based on status and kind.
- 0513f58: Adds background task reset commands which can be use to cancel one or more background tasks.

## 1.6.0

### Minor Changes

- 0e23aaf: Add deployment diagnostics commands

### Patch Changes

- 5a925ce: Remove default value for "reason" flag on access ensure command

## 1.5.0

### Minor Changes

- 9cc163a: Revert addition of the provisioning status for grants

## 1.4.1

### Patch Changes

- 18bc17a: Apply upgrade sdk to patched version which fixes an issue with the login command failing due to a nonce mismatch error.

## 1.4.0

### Minor Changes

- 2095990: Add support for supplying and displaying request reasons

## 1.3.1

### Patch Changes

- b71eeec: bump sdk version

## 1.3.0

### Minor Changes

- 96e8000: Add the manage > deployment > logs get/watch sub commands. Assists in filtering and tailing logs for a Common Fate deployment.
