docker rm -f $(docker ps -aq)
kill -9 $(lsof -t -i:4009)
for i in {5003..5066}
do
	kill -9 $(lsof -t -i:$i) &
	#echo $(lsof -t -i:$i)
done

rm default-ol/worker/worker.pid &
for i in {1..63}
do
	rm worker$i/worker/worker.pid &
done
