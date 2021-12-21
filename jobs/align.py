import pyspark.sql.functions as F
import pyspark.sql.types as T
from pyspark.sql import SparkSession
from pyspark import SparkConf
from subprocess import Popen, PIPE, STDOUT
import json


def wrapper(n: int = 8, overlap: int = 7, theta: int = 300, text: bool = False, cmd_path="./alignment"):
    return_text = {
        True: "--text=true",
        False: "--text=false"
    }[text]

    return_context = {
        True: "--context=true",
        False: "--context=false"
    }[text]

    def map_func(res, source, target, return_text=True):
        if return_text:
            return  list(map(lambda x: (
                int(x["source"].get("begin", 0)),
                int(x["source"].get("end", 0)),
                x["source"].get("text", ""),
                x["source"].get("before", ""),
                x["source"].get("after", ""),
                len(source),
                int(x["target"].get("begin", 0)),
                int(x["target"].get("end", 0)),
                x["target"].get("text", ""),
                x["target"].get("before", ""),
                x["target"].get("after", ""),
                len(target)
            ), res))
        else:
            return list(map(lambda x: (
                int(x["source"].get("begin", 0)),
                int(x["source"].get("end", 0)),
                int(x["target"].get("begin", 0)),
                int(x["target"].get("end", 0))
            ), res))

    def inner(source: str, target: str):
        msg = json.dumps({
            "source": source,
            "target": target,
            "config": json.dumps({
                "seeder": {"hash": {"n": n, "overlap": overlap}},
                "extender": {"range": {"theta": theta}}
            })
        })
        stdout = ""
        try:
            p = Popen([cmd_path, "read", return_text, return_context], stdout=PIPE, stdin=PIPE, stderr=STDOUT)
            stdout, _ = p.communicate(input=msg.encode("utf-8"))
            p.kill()
            stdout = json.loads(stdout.decode("utf-8"))
            return list(map_func(stdout, source, target))
        except Exception as e:
            return [(0, 0, stdout, "", "", 0, 0, 0, stdout, "", "", 0)]
    
    return inner


def align(source, target, wrapper):
    return wrapper(source, target)


def process(sc, path, align_cmd):
    # create candidate pair dataframe
    df = (
        # load parquet
        sc.read.parquet(path)
        .select(
            F.col("doi_a"),
            F.col("doi_b"),
            F.col("text_a"),
            F.col("text_b")
        )
    )
    
    df = (
        df
        .withColumn(
            "seeds",
            F.udf(
                lambda source, target: align(source, target, align_cmd),
                T.ArrayType(T.StructType([
                    T.StructField("begin_a", T.IntegerType(), False),
                    T.StructField("end_a", T.IntegerType(), False),
                    T.StructField("text_a", T.StringType(), True),
                    T.StructField("before_a", T.StringType(), True),
                    T.StructField("after_a", T.StringType(), True),
                    T.StructField("doc_length_a", T.IntegerType(), True),
                    T.StructField("begin_b", T.IntegerType(), False),
                    T.StructField("end_b", T.IntegerType(), False),
                    T.StructField("text_b", T.StringType(), True),
                    T.StructField("before_b", T.StringType(), True),
                    T.StructField("after_b", T.StringType(), True),
                    T.StructField("doc_length_b", T.IntegerType(), True),
                ]))
            )("text_a", "text_b")
        )
        .drop("text_a", "text_b")
    )

    df = (
        df
        .filter(F.size("seeds") > 0)
        .withColumn(
            "seeds",
            F.explode("seeds")
        )
        .select(
            F.col("doi_a"),
            F.col("doi_b"),
            F.col("seeds.*")
        )
    )
    return df


def run(sc, args):
    # args
    input_path = args[0]
    output_path = args[1]
    ngram_length = int(args[2])
    overlap = int(args[3])
    theta = int(args[4])
    # process
    cmd = wrapper(n=ngram_length, overlap=overlap, theta=theta, text=True, cmd_path="./jobs/alignment")

    # apply
    df = process(sc, input_path, cmd)
    df.show(vertical=True, truncate=True)
    df.write.mode("overwrite").parquet(output_path)


if __name__ == '__main__':
    # args
    INPUT = "stereo-joined.parquet/"
    OUTPUT = "stereo-aligned.parquet"
    BATCH = "00"
    SUBBATCH = ""
    NGRAM_LENGTH = 8
    NGRAM_OVERLAP = 7
    THETA = 250

    CMD = wrapper(n=NGRAM_LENGTH, overlap=NGRAM_OVERLAP, theta=THETA, text=True, cmd_path="./alignment")

    # spark session
    spark = (
        SparkSession
        .builder
        .config(conf=SparkConf())
        .getOrCreate()
    )
    # process
    input_batched = INPUT + BATCH + "/part-*"+SUBBATCH+"-*.parquet" 
    output_batched = OUTPUT + "/" + BATCH
    (
        process(spark, input_batched, CMD)
        .write
        .mode('overwrite')
        .parquet(output_batched)
    )
