import pyspark.sql.types as T
import pyspark.sql.functions as F
from pyspark.sql import SparkSession
from pyspark import SparkConf

import lxml.etree as ET
from urllib import parse

def extract(inp: str):
    try:
        inp = inp.split("</TEI>")[0]+"</TEI>"
        inp = bytes(inp, 'utf-8')
        tree = ET.fromstring(inp)
        for ref in tree.iter("{http://www.tei-c.org/ns/1.0}ref"):
            ref.text = ""
        text = []
        for elem in tree.iter("{http://www.tei-c.org/ns/1.0}p"):
            text.append(ET.tostring(elem, encoding="utf-8", method='text').decode("utf-8"))

        return " ".join(text)
    except:
        return ""


def process(sc, input_path):
    return (
        sc.createDataFrame(
            data=(
                sc.sparkContext
                # read file by line
                .textFile(input_path)
                # remove rows that are non-content lines from the original file
                .filter(lambda x: (x is not '{') and (x is not '}'))
                # separate DOI and file content
                .map(lambda x: (x.split(":")[0].strip(), ':'.join(x.split(":")[1:]).strip()))
                # remove leading and trailing " and remove hardcoded escapes
                .map(lambda x: (parse.unquote(x[0][1:-9]), x[1][1:-1].encode().decode('unicode_escape')))
                .map(lambda x: (x[0], extract(x[1])))
            ),
            schema=T.StructType([
                T.StructField("doi", T.StringType(), True),
                T.StructField("content", T.StringType(), True)
            ])
        )
        .where('content != ""')
        # collapse whitespace
        .withColumn(
            "content",
            F.regexp_replace("content", " +", " ")
        )
        # trim leading and trailing whitespace
        .withColumn(
            "content",
            F.trim(F.col("content"))
        )
    )


def run(sc, args):
    # Arguments
    input_path = args[0]
    output_path = args[1]
    df = process(sc, input_path)
    df.show()
    df.write.mode("overwrite").parquet(output_path)


if __name__ == '__main__':
    # args
    INPUT = ""
    OUTPUT = "stereo-grobid-preprocessed.parquet"
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
