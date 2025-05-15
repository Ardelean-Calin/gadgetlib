.PHONY: test test-integration test-all

test:
	go test ./gadget

test-integration:
	docker run -it --rm \
	  --cap-add=SYS_ADMIN \
	  --cap-add=SYS_MODULE \
	  --cap-add=NET_ADMIN \
	  --device /dev/fuse \
	  --security-opt apparmor=unconfined \
	  -v $(PWD):/usr/src/app \
	  -v /sys/kernel/config:/sys/kernel/config \
	  -w /usr/src/app \
	 	golang:1.24 /bin/bash -c "go test ./gadget -tags=integration"

	 test-all: test test-integration

