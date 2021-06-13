.PHONY: build-nextjs
build-nextjs:
	cd nextjs; \
	yarn install; \
	NEXT_TELEMETRY_DISABLED=1 yarn run export

.PHONY: build
build: build-nextjs
	go build .

.PHONY: build-nextjs-win
build-nextjs-win:
	cd nextjs && yarn install && npx next telemetry disable && yarn run export-win

.PHONY: build-win
build-win: build-nextjs-win
	go build .