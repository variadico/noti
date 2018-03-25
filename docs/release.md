# Noti release

This is the internal process I go through to release a version of Noti. I'm
just writing this down for myself.

## Tests

Make sure latest dev is green on CI.

* https://circleci.com/gh/variadico/noti
* https://ci.appveyor.com/project/variadico/noti

## Increment version

* docs/CHANGELOG.md
* docs/man/noti.1.md
* docs/man/noti.yaml.5.md

## Merge to master

```
git checkout master
git merge dev --ff-only
git push origin master
```

## Tests

Make sure latest master is green on CI.

* https://circleci.com/gh/variadico/noti
* https://ci.appveyor.com/project/variadico/noti

## Double check

Fix anything that might have broken like CI or URLs in docs. Last chance to
change anything.

## Tag release

Once everything is ready, tag the release.

```
git tag 1.2.3
git push origin 1.2.3
```

## Edit GitHub release information

* Click on Releases > 1.2.3 > Edit tag.
* Make the release title 1.2.3.
* Copy and paste the changes from CHANGELOG.md into the description box.

Create release tarballs.

```
make release
```

Upload files.

## Eventually update Homebrew

Read this: https://github.com/Homebrew/homebrew-core/blob/master/.github/CONTRIBUTING.md#submit-a-123-version-upgrade-for-the-foo-formula

And this: https://github.com/Homebrew/brew/blob/master/share/doc/homebrew/How-To-Open-a-Homebrew-Pull-Request-(and-get-it-merged).md#create-your-pull-request-from-a-new-branch
