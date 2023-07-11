yarn build

# /static - for build embed
copy-item -r -force ./build/* ../../static/html/admin

# /bin/static - for run, test binary
# copy-item -r -force ./build/* ../../bin/static/html/admin
