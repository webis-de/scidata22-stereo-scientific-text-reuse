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
  --name "stereo Finalize" \
  --master yarn \
  --deploy-mode cluster \
  --packages org.apache.hadoop:hadoop-aws:2.7.1 \
  --num-executors 100 \
  --executor-cores 4 \
  --executor-memory 16g \
  --driver-memory 32g \
  --conf spark.yarn.executor.memoryOverhead=8000 \
  --conf spark.yarn.driver.memoryOverhead=8000 \
  --archives venv.tar.gz#venv \
  --conf spark.sql.shuffle.partitions=5000 \
  --conf spark.yarn.appMasterEnv.PYSPARK_PYTHON=./venv/bin/python \
  --conf spark.network.timeout=600 \
  --py-files "jobs.zip" \
  jobs/finalize.py
