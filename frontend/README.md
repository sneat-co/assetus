# assetus frontend

Nx workspace for the assetus frontend: the standalone `assetus-app` and the
publishable `@sneat/extension-assetus` library.

- **Nx** 22 · **Angular** 21 · **Ionic** 8 · **pnpm**

## Setup

```bash
pnpm install
```

## Common tasks

```bash
pnpm exec nx serve assetus-app          # run the app locally
pnpm exec nx build ext-assetus          # build the publishable library
pnpm exec nx run-many -t lint test build
pnpm exec nx e2e assetus-app-e2e        # end-to-end tests
```

## Layout

```
frontend/
├── apps/
│   └── assetus-app/        # standalone assetus.app (Ionic shell)
└── libs/
    └── ext-assetus/        # @sneat/extension-assetus (publishable)
```

> Projects are generated incrementally during the extraction; see the repo
> root README for the overall plan.
