# Disabled workflows

`deploy-docker.yml` is deliberately stored outside `.github/workflows/` so GitHub Actions cannot run it.

It formerly published `multiversx/chainsimulator` images to Docker Hub on release events. NewArc policy permits workflows only on the `NewArc` branch and does not permit registry publication from this repository. Re-enable it only after a separate reviewed Xorewa registry and release policy is approved.
