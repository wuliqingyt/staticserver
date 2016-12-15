source ./config.sh
if [[ -z "$SERVER_NAME" ]]; then
    echo "Empty server name"
fi
kill `ps -ef | grep -w $SERVER_NAME | grep -v grep | awk '{print $2}'` 2>/dev/null
