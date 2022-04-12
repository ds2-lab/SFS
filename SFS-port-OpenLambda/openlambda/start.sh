./ol worker &
for i in {1..63}
do
	./ol worker --path=worker$i &
done
