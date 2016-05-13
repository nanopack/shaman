#!/usr/bin/env bash
set -e

# for versioning
getCurrCommit() {
  echo `git rev-parse HEAD | tr -d "[ \r\n\']"`
}

# for versioning
getCurrTag() {
  echo `git describe --always --tags --abbrev=0 | tr -d "[v\r\n]"`
}

# for versioning
getCurrBranch() {
  echo `git rev-parse --abbrev-ref HEAD | tr -d "[\r\n ]"`
}

# for versioning
commit=$(getCurrCommit)
branch=$(getCurrBranch)
tag=$(getCurrTag)

# try and use the correct MD5 lib (depending on user OS darwin/linux)
MD5=$(which md5 || which md5sum)

# build shaman
echo "Building SHAMAN and uploading it to 's3://tools.nanopack.io/shaman'"
gox -ldflags="-X main.version=${tag} -X main.branch=${branch} -X main.commit=${commit}" -osarch "linux/amd64" -output="./build/{{.OS}}/{{.Arch}}/shaman"

# look through each os/arch/file and generate an md5 for each
echo "Generating md5s..."
for os in $(ls ./build); do
  for arch in $(ls ./build/${os}); do
    for file in $(ls ./build/${os}/${arch}); do
      cat "./build/${os}/${arch}/${file}" | ${MD5} | awk '{print $1}' >> "./build/${os}/${arch}/${file}.md5"
    done
  done
done

# upload to AWS S3
echo "Uploading builds to S3..."
aws s3 sync ./build/ s3://tools.nanopack.io/shaman --grants read=uri=http://acs.amazonaws.com/groups/global/AllUsers --region us-east-1

#
echo "Cleaning up..."

# remove build
[ -e "./build" ] && \
  echo "Removing build files..." && \
  rm -rf "./build"
