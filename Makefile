.PHONY: test

test: test-units test-integration

test-units:
	go test ./gadget

test-integration:
	sudo modprobe configfs
	sudo modprobe dummy_hcd
	sudo modprobe libcomposite
	sudo go test ./gadget -tags=integration

# test-docker:
# 	docker run -it --rm \
# 	  --cap-add=SYS_ADMIN \
# 	  --cap-add=SYS_MODULE \
# 	  --cap-add=NET_ADMIN \
# 	  --device /dev/fuse \
# 	  --security-opt apparmor=unconfined \
# 	  -v $(PWD):/usr/src/app \
# 	  -v /sys/kernel/config:/sys/kernel/config \
# 	  -w /usr/src/app \
# 	 	golang:1.24 /bin/bash -c "go test ./gadget -tags=integration"

