import pyspark.sql.functions as F
import pyspark.sql.types as T
from pyspark.sql import SparkSession
from pyspark import SparkConf


def process(sc, input_path):
    # load data
    df = (
        # read file
        sc.read.parquet(input_path)
        # select columns
        .select(
            F.col("doi"),
            F.col("hash")
        )
        # reduce chunks to individual documents
        .rdd
        .map(lambda x: (x[0], [int(y[0]) for y in x[1]]))
        .reduceByKey(lambda a, b: a + b)
        .toDF(schema=T.StructType([
            T.StructField("doi", T.StringType(), False),
            T.StructField("hash", T.ArrayType(T.IntegerType()), False)
        ]))
    )
    return df


def run(sc, args):
    # args
    input_path = args[0]
    output_path = args[1]
    # process
    df = process(sc, input_path)
    df.show()
    df.write.mode("overwrite").parquet(output_path)


if __name__ == '__main__':
    # args
    INPUT = "stereo-hashed.parquet/*"
    OUTPUT = "stereo-reduced.parquet"
    # spark session
    spark = (
        SparkSession
        .builder
        .config(conf=SparkConf())
        .getOrCreate()
    )
    # process
    (
        process(spark, INPUT)
        .write
        .mode('overwrite')
        .parquet(OUTPUT)
    )
