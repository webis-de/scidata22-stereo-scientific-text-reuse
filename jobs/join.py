import pyspark.sql.functions as F
from pyspark.sql import SparkSession
from pyspark import SparkConf


def process(sc, pairs, texts):

    pairs = (
        sc.read.parquet(pairs)
        .select(
            F.col("A").alias("doi_a"),
            F.col("B").alias("doi_b")
        )
    )

    texts = (
            sc.read.parquet(texts)
            .persist()
    )

    return (
        pairs.alias("pairs")    
        .join(
            texts.alias("texts").select(F.col("doi"), F.col("content").alias("text_a")),
            F.col("pairs.doi_a") == F.col("texts.doi"),
            'left'
        )
        .drop("doi")
        .alias("pairs")
        .join(
            texts.alias("texts").select(F.col("doi"), F.col("content").alias("text_b")),
            F.col("pairs.doi_b") == F.col("texts.doi"),
            'left'
        )
        .drop("doi")
    )


def run(sc, args):
    # args
    input_pairs = args[0]
    input_texts = args[1]
    output_path = args[2]

    df = process(sc, input_pairs, input_texts)
    df.explain()
    df.write.mode("overwrite").parquet(output_path)


if __name__ == '__main__':
    # args
    PAIR_PATH = "stereo-paired.parquet"
    TEXT_PATH = "stereo-filtered.parquet/*"
    OUTPUT_PATH = "stereo-joined.parquet"
    BATCH = "00"
    NUM_PARTITIONS = 50_000

    spark = (
        SparkSession
        .builder
        .config(conf=SparkConf())
        .getOrCreate()
    )
    # process
    pair_batched = PAIR_PATH + "/part-*" + BATCH + "-*.parquet"
    output_batched = OUTPUT_PATH + "/" + BATCH
    (
        process(spark, pair_batched, TEXT_PATH)
        .write
        .mode('overwrite')
        .parquet(output_batched)
    )
