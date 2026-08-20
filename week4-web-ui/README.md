# week4-web-ui

Vue 3 + Vite dashboard for the PostgreSQL PaaS platform. It talks to the
`week3-paas-api` Go service and lets operators and registered users sign in,
provision CloudNativePG instances, inspect connection credentials and delete
instances. Each non-admin user only sees the instances they created.

## Requirements

- Node.js `^22.18.0 || >=24.12.0`
- A running `week3-paas-api` (on `http://localhost:8080` for local development)

## Setup

```sh
npm install
npm run dev
```

```sh
npm run build     # production bundle in dist/
npm run preview   # serve the production bundle
```

## How the API is reached

The app only ever calls relative paths such as `/api/v1/login`; it has no idea
where the backend runs. `localhost` must not appear in the bundle, because a
visitor's browser would resolve it to their own machine instead of the cluster.

- **Development:** the Vite dev server proxies `/api` to
  `http://localhost:8080` (see `vite.config.js`), so requests stay same-origin
  and CORS never comes into play.
- **Production:** serve `dist/` and route `/api` to the API Service from the
  same host, typically with an Ingress rule per path. Because it is one origin,
  the API's CORS allowlist is irrelevant there too.

## Configuration

| Variable | Default | Description |
| --- | --- | --- |
| `VITE_DEV_API_TARGET` | `http://localhost:8080` | Where `npm run dev` proxies `/api` |
| `VITE_API_BASE_URL` | _(empty)_ | Absolute API origin. Only set it to call the API cross-origin; that origin then needs to be in the API's CORS allowlist |

## API endpoints used

| Method | Path | Purpose |
| --- | --- | --- |
| `POST` | `/api/v1/register` | Create a platform account, returns `{ username }` |
| `POST` | `/api/v1/login` | Exchange credentials for a JWT (valid one hour) |
| `GET` | `/api/v1/instances` | List instances owned by the signed-in user (admins see all), returns `{ items, count }` |
| `POST` | `/api/v1/instances` | Create an instance, returns `202 Accepted` |
| `DELETE` | `/api/v1/instances/:name` | Start deletion, returns `202 Accepted` |
| `GET` | `/api/v1/instances/:name/connection` | Read credentials from the CNPG secret |

The JWT is stored in `localStorage`. The UI reads `sub` (username) and `role`
(`admin` or `user`) from the payload. Any `401` response clears the token
and returns the user to the login screen. After a successful register the
dashboard signs the new user in immediately.

How those calls reach Kubernetes, and how one user's instances stay
invisible to another, is in [../architecture.md](../architecture.md).

## Project structure

```
src/
├── api/client.js              # fetch wrapper, one method per API endpoint
├── composables/useToasts.js   # shared toast notification state
├── components/
│   ├── AppHeader.vue          # top bar with session info and logout
│   ├── LoginView.vue          # sign-in and create-account screen
│   ├── PlatformStats.vue      # instance/health/replica summary
│   ├── CreateInstanceForm.vue # create form with client-side validation
│   ├── InstanceList.vue       # searchable list, empty and loading states
│   ├── InstanceCard.vue       # single instance summary and actions
│   ├── ConnectionModal.vue    # credentials with reveal and copy
│   ├── ConfirmDeleteDialog.vue# type-to-confirm delete guard
│   └── ui/                    # BaseButton, BaseInput, BaseModal, AppCard,
│                              # StatusBadge, ToastStack
└── App.vue                    # state and API orchestration
```

Styling is Tailwind CSS v4 via `@tailwindcss/vite`; the only handwritten CSS
lives in `src/assets/main.css`.
