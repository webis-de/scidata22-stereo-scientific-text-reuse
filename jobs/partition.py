import pyspark.sql.functions as F
from pyspark.sql import SparkSession
from pyspark import SparkConf


def partitioner(n):
    def partitioner_(x):
        return x % n
    return partitioner_


def process(sc, input_path):
    return (
        sc.read.parquet(input_path)
        .select(
            F.col("doi"),
            F.col("hash")
        )
        .select(
            F.col("doi").alias("doi"),
            F.explode("hash").alias("hash")
        )
        .groupby("hash")
        .agg(F.collect_set("doi").alias("doi"))
        .select(
            F.col("hash"),
            F.col("doi")
        )
    )


def run(sc, args):
    # args
    input_path = args[0]
    output_path = args[1]
    # process
    df = process(sc, input_path)
    df.explain()
    df.write.mode('overwrite').parquet(output_path)


if __name__ == '__main__':
    # args
    INPUT = "stereo-reduced.parquet"
    OUTPUT = "stereo-partitioned.parquet"
    # spark session
    spark = (
        SparkSession
        .builder
        .config(conf=SparkConf())
        .getOrCreate()
    )
    (
        process(spark, INPUT)
        .write
        .mode('overwrite')
        .parquet(OUTPUT)
    )
