import pyspark.sql.functions as F
from pyspark.sql import SparkSession
from pyspark import SparkConf


def process(sc, input_path, doi_path, max_words):
    dois = (
        sc.read.text(doi_path)
        .select(
            F.col("value").alias("doi")
        )
    )
    return (
        dois
        .join(
            (
                # load parquet
                sc.read.parquet(input_path)
                # select columns
                .select(
                    F.col("doi"),
                    F.col("content")
                )
            ),
            on="doi", how="inner"
        )
        .filter(F.size(F.split(F.col("content"), " ")) < max_words)
        .select(
            F.col("doi"),
            F.col("content")
        )
    )


def run(sc, args):
    # args
    input_path = args[0]
    output_path = args[1]
    doi_path = args[2]
    max_words = args[3]
    # process
    (
        process(sc, input_path, doi_path, max_words)
        .write
        .mode('overwrite')
        .parquet(output_path)
    )


if __name__ == '__main__':
    # args
    INPUT = "stereo-grobid-preprocessed.parquet/*"
    OUTPUT = "grobid-filtered.parquet"
    DOI_PATH = "grobid_dois.txt"
    MAX_WORDS = 65_000
    # spark session
    spark = (
        SparkSession
        .builder
        .config(conf=SparkConf())
        .getOrCreate()
    )
    # process
    (
        process(spark, INPUT, DOI_PATH, MAX_WORDS)
        .write
        .mode('overwrite')
        .parquet(OUTPUT)
    )
