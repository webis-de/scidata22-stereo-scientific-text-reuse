# CLI targets
clean:
	$(MAKE) -C tools/alignment clean
clean:
	$(MAKE) -C tools/alignment test
build:
	$(MAKE) -C tools/alignment build

# Job Targets

preprocess-local:
	python3 main.py --job preprocess --job_args test_grobid test_preprocessed

preprocess-cluster:
	sh ./scripts/submit-preprocess.sh

filter-local:
	python3 main.py --job filter --job_args test_partitioned test_filtered remove.parquet

filter-cluster:
	sh ./scripts/submit-filter.sh

vectorize-local:
	python3 main.py --job vectorize --job_args test_preprocessed test_vectorized 50

vectorize-cluster:
	sh ./scripts/submit-vectorize.sh

hash-local:
	python3 main.py --job hash --job_args test_vectorized test_hashed 10

hash-cluster:
	sh ./scripts/submit-hash.sh

reduce-local:
	python3 main.py --job reduce --job_args test_hashed test_reduced

reduce-cluster:
	sh ./scripts/submit-reduce.sh

partition-local:
	python3 main.py --job partition --job_args test_reduced test_partitioned 10

partition-cluster:
	sh ./scripts/submit-partition.sh

pair-local:
	python3 main.py --job pair --job_args test_partitioned test_paired 10

pair-cluster:
	sh ./scripts/submit-pair.sh

join-local:
	python3 main.py --job join --job_args test_unified test_preprocessed test_joined

join-cluster:
	sh ./scripts/submit-join.sh

align-local: build
	python3 main.py --job align --job_args test_joined test_aligned 8 7 300

align-cluster: build
	sh ./scripts/submit-align.sh

metadata-local:
	python3 main.py --job metadata --job_args test_oag test_metadata

metadata-cluster:
	sh ./scripts/submit-metadata.sh

unify-local:
	python3 main.py --job unify --job_args test_aligned test_metadata

unify-cluster:
	sh ./scripts/submit-unify.sh

finalize-local:
	python3 main.py --job finalize --job_args test_unified

finalize-cluster:
	sh ./scripts/submit-finalize.sh
