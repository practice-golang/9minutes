yarn build

# /static - for build embed
remove-item -r -force ../../static/html/admin/_app
copy-item -r -force ./build/* ../../static/html/admin

# /bin/static - for run, test binary
# copy-item -r -force ./build/* ../../bin/static/html/admin
