all: simd-env

simd-env: simd-rmi
	docker build --tag cosmossdk/simd -f simd-env/Dockerfile \
		$(shell git rev-parse --show-toplevel) --build-arg GH_TOKEN="${GH_ACCESS_TOKEN}"

simd-dlv: simd-rmi
	docker build --tag cosmossdk/simd -f simd-dlv/Dockerfile \
		$(shell git rev-parse --show-toplevel) --build-arg GH_TOKEN="${GH_ACCESS_TOKEN}"

simd-rmi:
	docker rmi cosmossdk/simd 2>/dev/null; true

.PHONY: all simd-env simd-dlv simd-rmi
