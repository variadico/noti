# Noti release

This is the internal process I go through to release a version of Noti. I'm just
writing this down for myself.

## Run tests

```
cd variadico/noti
./lint.sh
./cmd/noti/test.sh
go test ./...
```

## Increment version

* CHANGELOG.md
* noti.go
* README.md
* gen_bin.sh

## Merge to mainline

```
git checkout master
git merge dev --ff-only

# Repeat Run Tests. Make sure it still works.
go install github.com/variadico/noti/cmd/noti

# Go check make sure you didn't screw anything up. Make sure URLs in doc resolve
# correctly.
git push origin master

git tag v1.2.3
git push origin v1.2.3
```

## Edit GitHub release information

* Click on Releases > v1.2.3 > Edit tag.
* Make the release title v1.2.3.
* Copy and paste the changes from CHANGELOG.md into the description box.

Upload binaries.

```
cd $GOPATH/bin
tar -czf noti1.2.3.darwin-amd64.tar.gz noti
mv noti1.2.3.darwin-amd64.tar.gz ~/Desktop

docker run --rm -it -v $GOPATH:/go golang:1.6.2 /bin/bash
cd $GOPATH/bin
rm noti

go install github.com/variadico/noti/cmd/noti
# Make sure it works.
noti -h

tar -czf noti1.2.3.linux-amd64.tar.gz noti
mv noti1.2.3.linux-amd64.tar.gz ~/Desktop
```

## Eventually update Homebrew

Read this: https://github.com/Homebrew/homebrew-core/blob/master/.github/CONTRIBUTING.md#submit-a-123-version-upgrade-for-the-foo-formula

And this: https://github.com/Homebrew/brew/blob/master/share/doc/homebrew/How-To-Open-a-Homebrew-Pull-Request-(and-get-it-merged).md#create-your-pull-request-from-a-new-branch
