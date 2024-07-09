#!/bin/bash

# 设置服务器的IP地址和用户名
SERVER_IP="129.146.130.68"
USERNAME="ubuntu"

DB_USER="root"
DB_PASS="123456"
DB_NAME="aion"

# 第一步，编译后端代码
echo "Compiling backend code..."
go build -o aion main.go
if [ $? -ne 0 ]; then
    echo "Backend compilation failed!"
    exit 1
fi

# 第二步，编译前端代码
echo "Compiling frontend code..."
export NODE_OPTIONS=--openssl-legacy-provider
cd frontend && npm run build
if [ $? -ne 0 ]; then
    echo "Frontend compilation failed!"
    exit 1
fi

# 打包前端代码
echo "Packaging frontend code..."
tar -czf dist.tgz dist/
if [ $? -ne 0 ]; then
    echo "Packaging frontend code failed!"
    exit 1
fi

cd .. || exit

echo "Exporting database..."
mysqldump -u$DB_USER -p$DB_PASS $DB_NAME > aion.sql
if [ $? -ne 0 ]; then
    echo "Database export failed!"
    exit 1
fi

tar -czf aion.sql.tgz aion.sql
if [ $? -ne 0 ]; then
    echo "tar sql file failed!"
    exit 1
fi

# 第三步，上传代码到服务器
echo "Uploading code to server..."
scp aion app.ini aion.sql.tgz frontend/dist.tgz ${USERNAME}@${SERVER_IP}:~
if [ $? -ne 0 ]; then
    echo "Uploading code failed!"
    exit 1
fi

# 第四步，解压前端代码到/opt/dist目录
echo "Deploying frontend code on server..."
ssh ${USERNAME}@${SERVER_IP} << 'EOF'
    tar -xzf dist.tgz
    tar -xzf aion.sql.tgz
    sudo rm -rf /opt/dist
    sudo mv dist /opt/dist
    sudo mv aion app.ini aion.sql /opt/
EOF
if [ $? -ne 0 ]; then
    echo "Deploying frontend code failed!"
    exit 1
fi

echo "Importing database on server..."
ssh ${USERNAME}@${SERVER_IP} << EOF
    mysql -u$DB_USER -p$DB_PASS $DB_NAME < /opt/aion.sql
EOF
if [ $? -ne 0 ]; then
    echo "Database import failed!"
    exit 1
fi

# 第五步，重启服务
echo "Restarting service on server..."
ssh ${USERNAME}@${SERVER_IP} /home/ubuntu/aion.sh restart
if [ $? -ne 0 ]; then
    echo "Service restart failed!"
    exit 1
fi

echo "Deployment successful!"

# 清理
rm aion frontend/dist.tgz aion.sql aion.sql.tgz