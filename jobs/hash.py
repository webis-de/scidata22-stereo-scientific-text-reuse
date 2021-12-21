import pyspark.sql.functions as F
import pyspark.sql.types as T
from pyspark.sql import SparkSession
from pyspark import SparkConf
from pyspark.ml.feature import MinHashLSH
import numpy as np


def process(sc, input_path, num_hashes):
    # create dataframe
    df = (
        # load parquet
        sc.read.parquet(input_path)
        # rename columns to model naming scheme
        .select(
            F.col("doi"),
            F.col("index"),
            F.col("features")
        )
        # filter empty rows
        .filter(
            F.udf(
                lambda x: bool(np.sum(x.toArray()) != 0),
                T.BooleanType()
            )("features")
        )
    )

    return (
        MinHashLSH(inputCol="features", outputCol="hash", numHashTables=num_hashes)
        # fit minhash model on text chunks
        .fit(df)
        # calculate hashes for each chunk
        .transform(df)
        # drop original features
        .drop("features")
    )


def run(sc, args):
    # args
    input_path = args[0]
    output_path = args[1]
    num_hashes = int(args[2])
    # process
    df = process(sc, input_path, num_hashes)
    df.show()
    df.write.mode("overwrite").parquet(output_path)


if __name__ == '__main__':
    # args
    INPUT = "stereo-vectorized.parquet/*"
    OUTPUT = "stereo-hashed.parquet"
    NUM_HASHES = 5

    # spark session
    spark = (
        SparkSession
        .builder
        .config(conf=SparkConf())
        .getOrCreate()
    )
    # process
    (
        process(spark, INPUT, NUM_HASHES)
        .write
        .mode('overwrite')
        .parquet(OUTPUT)
    )
