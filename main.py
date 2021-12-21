from pyspark.sql import SparkSession
from pyspark import SparkConf
import argparse
import importlib


parser = argparse.ArgumentParser()
parser.add_argument('--job', type=str, required=True)
parser.add_argument('--job_args', nargs='*')
args = parser.parse_args()

spark = (
    SparkSession
    .builder
    .master("local")
    .config(conf=SparkConf())
    .getOrCreate()
)

job_module = importlib.import_module('jobs.%s' % args.job)
job_module.run(spark, args.job_args)
