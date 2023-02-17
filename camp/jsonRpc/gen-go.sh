#! /bin/bash


#python3 -m grpc_tools.protoc --python_out=. --grpc_python_out=. -I. lot_service.proto
temp=$(pwd)
array=(${temp//// })
unset 'array[${#array[@]}-1]'
unset 'array[${#array[@]}-1]'
uploadDir='/'
for var in ${array[@]};do
  uploadDir+=$var'/'
done

echo $uploadDir

unset 'array[${#array[@]}-1]'
baseDir='/'
for var in ${array[@]};do
  baseDir+=$var'/'
done
echo $baseDir
cd ..
cp -r jsonRpc $uploadDir
cd $uploadDir''jsonRpc


protoc --go_out=plugins=grpc:.  *.proto






ignore=("gen-go.sh")
sub=$(ls *.proto)
temp=(${sub[@]} ${ignore[*]})
for element in ${temp[@]}; do
    rm $element
done
