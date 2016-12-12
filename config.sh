SERVER_NAME=$(grep 'server_name=' Makefile | awk -F= '{ print $2  }')
