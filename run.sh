source ./config.sh

sh stop.sh
while true
do
    if [[ "$(sh ll.sh | wc -l)" -eq 0 ]]; then
        break
    fi
    usleep 2000
done

nohup ./$SERVER_NAME &
