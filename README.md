# mesWin

Photo-based window measurement for Rehau Synego uPVC systems. Field technicians capture inner/outer photos with a calibration target; Go + ONNX CV extracts dimensions to ±5 mm; exports feed sales and manufacturing.

## Stack

| Layer | Tech |
|-------|------|
| Frontend | SvelteKit + Tailwind + Capacitor |
| Backend | Go + PocketBase (embedded) |
| CV | ONNX YOLO in Go (stub active; training in `cv/`) |

## Quick start

### Backend

```bash
cd backend
go run . serve --http=127.0.0.1:8090
```

Admin UI: http://127.0.0.1:8090/_/

### Frontend

```bash
cd frontend
cp .env .env.local   # optional
npm install
npm run dev
```

App: http://localhost:5173

### First-time setup (PocketBase admin)

1. Open http://127.0.0.1:8090/_/ and create admin account.
2. Create an **organization** record.
3. Create a **user** (Auth) with `organization` relation and `role` = `field`.
4. Sign in via the frontend `/login`.

## Project layout

```
backend/          Go + PocketBase binary, custom API, CV stub
frontend/         SvelteKit app (web + Capacitor-ready)
deploy/           Docker image + compose
cv/               Model training (future)
docs/             vision.md, taskflow.md, skills.md, agents.md
```

## API (custom)

- `GET  /api/meswin/health`
- `POST /api/meswin/windows/:id/process` — run CV (auth required)
- `GET  /api/meswin/orders/:id/export.json` — ERP export (auth required)

## Docker

```bash
docker compose -f deploy/docker-compose.yml up --build
```

## Docs

- [vision.md](./vision.md) — goals & requirements
- [taskflow.md](./taskflow.md) — phased delivery
- [skills.md](./skills.md) — AI/dev domain skills
- [agents.md](./agents.md) — agent roles

## Accuracy target

±5 mm with printed 90° calibration triangle (100–150 mm, 10 mm red markings). 10 mm hard tolerance.
