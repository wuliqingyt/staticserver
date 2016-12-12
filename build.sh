source ./config.sh
make
if [[ $? -eq 0 ]]; then
    ./$SERVER_NAME
fi
