import pyspark.sql.functions as F
import pyspark.sql.types as T
from pyspark.sql import SparkSession
from pyspark import SparkConf
from pyspark.ml.feature import Tokenizer, StopWordsRemover, HashingTF


def sliding_window_curried(n, overlap):
    def sliding_window(iterable):
        inp = iter(iterable)
        win = []
        for e in range(0, n):
            try:
                win.append(next(inp))
            except StopIteration:
                break
        yield win
        for i, e in enumerate(inp, start=1):
            win = win[1:] + [e]
            if i % overlap == 0:
                yield win
    return lambda x: list(sliding_window(x))


def process(sc, input_path, ngram_length, num_features):
    # load data
    df = sc.read.parquet(input_path)
    df = (
        df
        # drop special characters
        .withColumn(
            "content",
            F.regexp_replace(F.col("content"), r"\\t|\\r|\\n|[.,;:{}=()]+", ' ')
        )
    )
    # tokenize input text
    df = (
        Tokenizer(inputCol='content', outputCol="tokens")
        .transform(df)
        .drop("content")
    )
    # remove stopwords
    df = (
        StopWordsRemover(inputCol="tokens", outputCol="filtered_tokens")
        .transform(df)
        .drop("tokens")
        # filter single letters
        .withColumn(
            "filtered_tokens",
            F.udf(
                lambda x: list(filter(None, [y if len(y) > 2 else None for y in x])),
                T.ArrayType(T.StringType())
            )("filtered_tokens")
        )
        # filter numbers
        .withColumn(
            "filtered_tokens",
            F.udf(
                lambda x: list(filter(None, [y if y.isnumeric() else None for y in x])),
                T.ArrayType(T.StringType())
            )("filtered_tokens")
        )
    )
    # create ngrams
    df = (
        df
        .withColumn(
            "ngrams",
            F.udf(sliding_window_curried(ngram_length, ngram_length), T.ArrayType(T.ArrayType(T.StringType())))("filtered_tokens")
        )
        .drop("filtered_tokens")
    )
    # explode ngram list to separate (numbered) rows
    df = df.select(
        "doi",
        F.posexplode("ngrams").alias("index", "ngrams")
    )
    # create chunk vectors
    df = (
        HashingTF(inputCol="ngrams", outputCol="features", numFeatures=num_features, binary=True)
        .transform(df)
        .drop("ngrams")
    )

    return df


def run(sc, args):
    # args
    input_path = args[0]
    output_path = args[1]
    ngram_length = int(args[2])
    num_features = 2**18
    # process
    df = process(sc, input_path, ngram_length, num_features)
    df.show()
    df.write.mode("overwrite").parquet(output_path)


if __name__ == '__main__':
    # args
    INPUT = "stereo-filtered.parquet/*"
    OUTPUT = "stereo-vectorized.parquet"
    NGRAM_LEN = 50
    NUM_FEATURES = 2**18

    # spark session
    spark = (
        SparkSession
        .builder
        .config(conf=SparkConf())
        .getOrCreate()
    )
    # process
    (
        process(spark, INPUT, NGRAM_LEN, NUM_FEATURES)
        .write
        .mode('overwrite')
        .parquet(OUTPUT)
    )
