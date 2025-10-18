# Dependency Management

This repository uses Renovate to keep dependency graphs current across the different services. Renovate is configured in `renovate.json` at the repository root and covers:

- Go modules under `services/manager/backend`
- npm workspaces under `services/manager/frontend` and `services/manager/openapi`

Pull requests opened by Renovate are labeled `dependencies`. A dependency dashboard issue is also created automatically, which you can use to trigger manual rechecks or to keep an eye on pending updates.

If you need to adjust the update cadence or add new packages to the manager-specific groups, edit `renovate.json` and push the change to the default branch. Consult the [Renovate documentation](https://docs.renovatebot.com/) for the full list of available options.
