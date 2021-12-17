import pyspark.sql.functions as F
import pyspark.sql.types as T
from pyspark.sql import SparkSession
from pyspark import SparkConf


def process(sc, input_cases, input_metadata, num_partitions):

    df_cases = (
        sc.read.parquet(input_cases)
        .select(
            F.col("doi_a"),
            F.col("doi_b"),
            F.col("begin_a"),
            F.col("end_a"),
            F.col("text_a"),
            F.col("before_a"),
            F.col("after_a"),
            F.col("doc_length_a"),
            F.col("begin_b"),
            F.col("end_b"),
            F.col("text_b"),
            F.col("before_b"),
            F.col("after_b"),
            F.col("doc_length_b"),
        )
    )

    df_metadata = (
        sc.read.parquet(input_metadata)
        .select(
            F.col("doi"),
            F.col("year"),
            F.col("board"),
            F.col("area"),
            F.col("discipline")
        )
    )

    df = (
        df_cases.alias("cases")
        .join(
            df_metadata.select(
                F.col("doi"),
                F.col("year").alias("year_a"),
                F.col("board").alias("board_a"),
                F.col("area").alias("area_a"),
                F.col("discipline").alias("discipline_a")
            ).alias("metadata"),
            F.col("cases.doi_a") == F.col("metadata.doi"),
            'left'
        )
        .drop("doi")
        .join(
            df_metadata.select(
                F.col("doi"),
                F.col("year").alias("year_b"),
                F.col("board").alias("board_b"),
                F.col("area").alias("area_b"),
                F.col("discipline").alias("discipline_b")
            ).alias("metadata"),
            F.col("cases.doi_b") == F.col("metadata.doi"),
            'left'
        )
        .drop("doi")
    )

    df = (
        df
        .select(
            F.col("doi_a"),
            F.col("begin_a"),
            F.col("end_a"),
            F.col("text_a"),
            F.col("before_a"),
            F.col("after_a"),
            F.col("doc_length_a"),
            F.col("year_a"),
            F.col("board_a").alias("field_a"),
            F.col("area_a"),
            F.col("discipline_a"),
            F.col("doi_b"),
            F.col("begin_b"),
            F.col("end_b"),
            F.col("text_b"),
            F.col("before_b"),
            F.col("after_b"),
            F.col("doc_length_b"),
            F.col("year_b"),
            F.col("board_b").alias("field_b"),
            F.col("area_b"),
            F.col("discipline_b"),
        )
    )

    return df.coalesce(num_partitions)


def run(sc, args): 
    case_path = args[0]
    metadata_path = args[1]
    process(sc, case_path, metadata_path, 1)


if __name__ == '__main__':
    # spark session
    spark = (
        SparkSession
        .builder
        .config(conf=SparkConf())
        .getOrCreate()
    )
    # argument
    INPUT_CASES = "stereo-aligned.parquet/*/*"
    INPUT_METADATA = "stereo-oag.parquet/*"
    NUM_PARTITIONS = 50
    OUTPUT = "stereo-corpus-final.jsonl"
    # process
    (
        process(spark, INPUT_CASES, INPUT_METADATA, NUM_PARTITIONS)
        .write
        .mode("overwrite")
        .option("compression", "gzip")
        .json(OUTPUT)
    )
