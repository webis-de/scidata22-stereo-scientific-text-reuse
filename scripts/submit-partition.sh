set -e
set -x

rm -f venv.tar.gz
rm -f jobs.zip

zip jobs.zip jobs

python3.7 -m venv venv
source venv/bin/activate
pip install -U pip
pip install -r requirements.txt
venv-pack -o venv.tar.gz

PYSPARK_PYTHON="./venv/bin/python" HADOOP_CONF_DIR=./hadoop-conf SPARK_HOME=./local-spark ./local-spark/bin/spark-submit \
  --conf spark.yarn.submit.waitAppCompletion=false \
  --name "stereo Partitioning" \
  --master yarn \
  --deploy-mode cluster \
  --packages org.apache.hadoop:hadoop-aws:2.7.1 \
   --num-executors 100 \
  --executor-cores 2 \
  --executor-memory 32g \
  --driver-memory 32g \
  --conf spark.yarn.executor.memoryOverhead=10000 \
  --conf spark.yarn.driver.memoryOverhead=10000 \
  --archives venv.tar.gz#venv \
  --conf spark.yarn.appMasterEnv.PYSPARK_PYTHON=./venv/bin/python \
  --conf spark.network.timeout=600 \
  --py-files "jobs.zip" \
  jobs/partition.py
