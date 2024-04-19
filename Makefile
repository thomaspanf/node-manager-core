beacon/ssz_types/types_encoding.go: beacon/ssz_types/types.go
	rm -f beacon/ssz_types/types_encoding.go
	sszgen --path ./beacon/ssz_types

.PHONY: clean
clean:
	rm -f beacon/ssz_types/types_encoding.go
