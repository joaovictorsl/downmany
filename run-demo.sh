# Build
make build
# Create datasets
./make_dataset.sh 10 ./dataset1
./make_dataset.sh 10 ./dataset2
./make_dataset.sh 10 ./dataset3
./make_dataset.sh 10 ./dataset4
echo "d" > ./dataset1/target
echo "a" > ./dataset2/target
echo "b" > ./dataset3/target
echo "c" > ./dataset4/target
# Run server
./bin/downmany -server -port 8000 > log_server 2>&1 &
pid_server=$!
# Run clients
./bin/downmany -file_hash 107 -server_addr localhost:8000 \
    -dataset ./dataset1 -port 8001 > log1 2>&1 &
pid_c1=$!

./make_dataset.sh 10 dataset2
./bin/downmany -file_hash 108 -server_addr localhost:8000 \
    -dataset ./dataset2 -port 8002 > log2 2>&1 &
pid_c2=$!

./make_dataset.sh 10 dataset3
./bin/downmany -file_hash 109 -server_addr localhost:8000 \
    -dataset ./dataset3 -port 8003 > log3 2>&1 &
pid_c3=$!

./make_dataset.sh 10 dataset4
./bin/downmany -file_hash 110 -server_addr localhost:8000 \
    -dataset ./dataset4 -port 8004 > log4 2>&1 &
pid_c4=$!
# Cleanup
rm -rf dataset*
