import pyspark.sql.functions as F
import pyspark.sql.types as T
from pyspark.sql import SparkSession
from pyspark import SparkConf

import itertools


def partitioner(n):
    def partitioner_(x):
        return x % n
    return partitioner_


def emit_pairs(group):
    for x, y in itertools.combinations(group, 2):
        if x > y:
            yield x, y
        else:
            yield y, x


def process(sc, input_path, num_partitions):
    df = (
        sc.read.parquet(input_path)
        .select(
            F.col("hash"),
            F.col("doi")
        )
        .filter(F.size("doi")> 1)
        .filter(F.size("doi") < 10_000)
    )
    return (
        df
        .rdd
        .keyBy(lambda x: x.hash)
        .partitionBy(numPartitions=num_partitions, partitionFunc=partitioner(num_partitions))
        .flatMap(lambda x: emit_pairs(x[1].doi))
        .toDF(schema=T.StructType([
            T.StructField("A", T.StringType(), False),
            T.StructField("B", T.StringType(), False)
        ]))
    )


def run(sc, args):
    # args
    input_path = args[0]
    output_path = args[1]
    num_partitions = 10
    # process
    df = process(sc, input_path, num_partitions)
    df.explain()
    (
        df
        .write
        .mode('overwrite')
        .parquet(output_path)
    )


if __name__ == '__main__':
    # args
    INPUT = "stereo-partitioned.parquet"
    OUTPUT = "stereo-paired.parquet"
    NUM_PARTITIONS = 40_000
    spark = (
        SparkSession
        .builder
        .config(conf=SparkConf())
        .getOrCreate()
    )
    # process
    (
        process(spark, INPUT, NUM_PARTITIONS)
        .write
        .mode('overwrite')
        .parquet(OUTPUT)
    )
