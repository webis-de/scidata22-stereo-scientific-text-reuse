import pyspark.sql.functions as F
import pyspark.sql.types as T
from pyspark.sql import SparkSession
from pyspark import SparkConf

import uuid


def process(sc, input_cases, input_texts, input_metadata, num_partitions):

    df_cases = (
        sc.read.json(input_cases)
        .select(
            F.col("doi_a"),
            F.col("begin_a"),
            F.col("end_a"),
            F.col("text_a"),
            F.col("before_a"),
            F.col("after_a"),
            F.col("doc_length_a"),
            F.col("year_a"),
            F.col("field_a"),
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
            F.col("field_b"),
            F.col("area_b"),
            F.col("discipline_b")     
        )
        .dropDuplicates()
    )

    df_cases = (
        df_cases.withColumn(
            "id",
            F.udf(lambda x: str(uuid.uuid4()), T.StringType())(F.col("doi_a"))
        )
        .select(
            F.col("id"),
            F.col("doi_a"),
            F.col("begin_a"),
            F.col("end_a"),
            F.col("text_a"),
            F.col("before_a"),
            F.col("after_a"),
            F.col("doc_length_a"),
            F.col("year_a"),
            F.col("field_a"),
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
            F.col("field_b"),
            F.col("area_b"),
            F.col("discipline_b"),
        )
    )

    df_cases = df_cases.coalesce(num_partitions)

    df_texts = (
        sc.read.parquet(input_texts)
        .select(
            F.col("doi").alias("doi_text"),
            F.col("content").alias("text")
        )
    )

    df_metadata = (
        sc.read.parquet(input_metadata)
        .select(
            F.col("doi"),
            F.col("year"),
            F.col("field"),
            F.col("area"),
            F.col("discipline")
        )
    )

    df_publications = (
        df_texts.alias("texts")
        .join(
            df_metadata.alias("metadata"),
            F.col("texts.doi_text") == F.col("metadata.doi"),
            'left'
        )
        .drop("doi_text")
        .dropDuplicates(subset=["doi"])
    )

    df_publications = df_publications.coalesce(num_partitions)

    return (
        df_cases,
        df_cases.drop("text_a", "text_b", "before_a", "before_b", "after_a", "after_b"),
        df_publications,
        df_publications.drop("text")
    )

def run(sc, args): 
    input_path = args[0]
    process(sc, input_path, 1)


if __name__ == '__main__':
    # spark session
    spark = (
        SparkSession
        .builder
        .config(conf=SparkConf())
        .getOrCreate()
    )
    
    # arguments
    INPUT_CASES = "stereo-corpus.jsonl"
    INPUT_TEXTS = "stereo-filtered.parquet"
    INPUT_METADATA = "stereo-oag.parquet"
    NUM_PARTITIONS = 50

    OUTPUT_CASES_FULL = "webis-stereo21/cases-full"
    OUTPUT_CASES_METADATA_ONLY = "webis-stereo21/cases-metadata"
    OUTPUT_PUBLICATIONS_FULL = "webis-stereo21/publications-full"
    OUTPUT_PUBLICATIONS_METADATA_ONLY = "webis-stereo21/publications-metadata-only"
    
    # process
    cases, cases_metadata_only, publications, publications_metadata_only = process(spark, INPUT_CASES, INPUT_TEXTS, INPUT_METADATA, NUM_PARTITIONS)
    (
        cases
        .write
        .mode("overwrite")
        .option("compression", "gzip")
        .json(OUTPUT_CASES_FULL)
    )
    
    (
        cases_metadata_only
        .write
        .mode("overwrite")
        .option("compression", "gzip")
        .json(OUTPUT_CASES_METADATA_ONLY)
    )
    
    (
        publications
        .write
        .mode("overwrite")
        .option("compression", "gzip")
        .json(OUTPUT_PUBLICATIONS_FULL)
    )

    (
        publications_metadata_only
        .write
        .mode("overwrite")
        .option("compression", "gzip")
        .json(OUTPUT_PUBLICATIONS_METADATA_ONLY)
    )
