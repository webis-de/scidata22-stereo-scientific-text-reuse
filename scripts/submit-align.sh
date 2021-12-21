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
  --name "stereo Aligning (1/100)" \
  --master yarn \
  --deploy-mode cluster \
  --packages org.apache.hadoop:hadoop-aws:2.7.1 \
  --num-executors 50 \
  --executor-cores 8 \
  --executor-memory 32g \
  --driver-memory 32g \
  --conf spark.yarn.executor.memoryOverhead=8000 \
  --conf spark.yarn.driver.memoryOverhead=10000 \
  --archives venv.tar.gz#venv \
  --conf spark.yarn.appMasterEnv.PYSPARK_PYTHON=./venv/bin/python \
  --conf spark.network.timeout=600 \
  --conf spark.sql.shuffle.partitions=5000 \
  --py-files "jobs.zip" \
  --files jobs/alignment#alignment \
  jobs/align.py
