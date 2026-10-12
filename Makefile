.PHONY: backend frontend run test

backend:
	cd backend && go run ./cmd/server

frontend:
	cd frontend && npm run dev

run:
	$(MAKE) -j2 backend frontend

test:
	cd backend && go test ./...
